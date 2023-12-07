package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"

	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func getVolume() int {
	cmd := exec.Command(`osascript`, `-e`, `output volume of (get volume settings)`)
	output, err := cmd.Output()
	if err != nil {
		return -1
	}
	volume := 0
	if _, err := fmt.Sscanln(string(output), &volume); err != nil {
		return -1
	}
	return volume
}

func setVolume(volume int) error {
	cmd := exec.Command(`osascript`, `-e`, `set volume output volume `+fmt.Sprint(volume))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func getMute() bool {
	cmd := exec.Command(`osascript`, `-e`, `output muted of (get volume settings)`)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	muted := false
	if _, err := fmt.Sscanln(string(output), &muted); err != nil {
		return false
	}
	return muted
}

func setMute(b bool) error {
	s := "without"
	if b {
		s = "with"
	}
	cmd := exec.Command(`osascript`, `-e`, fmt.Sprintf(`set volume %s output muted`, s))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func lockScreen() error {
	in := `
tell application "System Events"
	keystroke "q" using {control down, command down}
end tell`
	cmd := exec.Command(`osascript`)
	cmd.Stdin = strings.NewReader(in)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func observeScreenLock(b chan<- bool) error {
	cmd := exec.Command(`./lock`)
	r, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		br := bufio.NewScanner(r)
		for br.Scan() {
			line := br.Text()
			switch {
			case strings.Contains(line, `Screen Locked`):
				b <- true
			case strings.Contains(line, `Screen Unlocked`):
				b <- false
			}
			log.Println(br.Text())
		}
		log.Println(br.Err())
	}()

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	var accessories []*accessory.A

	{
		a := accessory.NewLightbulb(accessory.Info{
			Name: "MacOS Volume",
		})
		a.Lightbulb.On.OnValueRemoteUpdate(func(v bool) {
			setMute(!v)
		})
		ch := characteristic.NewBrightness()
		ch.OnValueRemoteUpdate(func(v int) {
			setVolume(v)
		})
		go func() {
			for range time.NewTicker(time.Second * 10).C {
				if v := getVolume(); v != -1 {
					ch.SetValue(v)
				}
				a.Lightbulb.On.SetValue(!getMute())
			}
		}()
		a.Lightbulb.AddC(ch.C)

		accessories = append(accessories, a.A)
	}

	{
		screenLocked := atomic.Bool{}
		a := accessory.NewSwitch(accessory.Info{
			Name: `MacOS Lock Screen`,
		})
		a.Switch.On.OnValueRemoteUpdate(func(v bool) {
			if !v {
				lockScreen()
			}
			a.Switch.On.SetValue(screenLocked.Load())
		})
		a.Switch.On.SetValue(true)

		b := make(chan bool)
		go func() {
			if err := observeScreenLock(b); err != nil {
				log.Println(err)
			}
		}()
		go func() {
			for locked := range b {
				a.Switch.On.SetValue(!locked)
				log.Println("屏幕锁定：", locked)
				screenLocked.Store(locked)
			}
		}()

		accessories = append(accessories, a.A)
	}

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	b := accessory.NewBridge(accessory.Info{
		Name: `MacOS HomeKit`,
	})

	// Create the hap server.
	server, err := hap.NewServer(fs, b.A, accessories...)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	server.Addr = ":6242" // MAHA
	log.Fatalln(server.ListenAndServe(ctx))
}
