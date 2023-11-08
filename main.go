package main

import (
	"fmt"
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

func main() {
	// Create the switch accessory.
	a := accessory.NewLightbulb(accessory.Info{
		Name: "MacBook Pro",
	})
	a.Lightbulb.On.OnValueRemoteUpdate(func(v bool) {
		setMute(!v)
	})
	ch := characteristic.NewBrightness()
	ch.OnValueRemoteUpdate(func(v int) {
		setVolume(v % 100)
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

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	// Create the hap server.
	server, err := hap.NewServer(fs, a.A)
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
	server.ListenAndServe(ctx)
}
