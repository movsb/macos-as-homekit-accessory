.PHONY: all
all: lock maha

lock: lock.swift
	swiftc $^

.PHONY: maha
maha:
	GOOS=darwin GOARCH=amd64 go build -o maha-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o maha-darwin-arm64
