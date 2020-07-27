#-include .env

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## run: run from code.
run:
	go run cmd/distribyted/main.go examples/conf_example.yaml

## build: build binary.
build: go-generate go-build

## test: execute all tests.
test:
	go test -v ./...

compile: go-generate go-compile

go-build:
	@echo "  >  Building binary..."
	go build $(LDFLAGS) -o bin/distribyted -tags "release" cmd/distribyted/main.go

go-generate:
	@echo "  >  Generating code files..."
	go generate ./...

go-compile:
	@echo "  >  Compiling for several platforms..."
	# 32-Bit Systems
	# FreeBDS
	#GOOS=freebsd GOARCH=386 go build $(LDFLAGS) -o bin/main-freebsd-386 -tags "release" cmd/distribyted/main.go
	# MacOS
	#GOOS=darwin GOARCH=386 go build $(LDFLAGS) -o bin/main-darwin-386 -tags "release" cmd/distribyted/main.go
	# Linux
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o bin/main-linux-386 -tags "release" cmd/distribyted/main.go
    # 64-Bit
	# FreeBDS
	#GOOS=freebsd GOARCH=amd64 go build $(LDFLAGS) -o bin/main-freebsd-amd64 -tags "release" cmd/distribyted/main.go
	# MacOS
	#GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/main-darwin-amd64 -tags "release" cmd/distribyted/main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/main-linux-amd64 -tags "release" cmd/distribyted/main.go
	
	# ARM
	GOOS=linux GOARCH=arm GOARM=5 go build $(LDFLAGS) -o bin/main-linux-arm -tags "release" cmd/distribyted/main.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo