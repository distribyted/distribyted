#-include .env

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
BIN_OUTPUT ?= bin/distribyted-$(VERSION)-`go env GOOS`-`go env GOARCH``go env GOEXE`
PROJECTNAME := $(shell basename "$(PWD)")

# Use linker flags to provide version/build settings
LDFLAGS=-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -linkmode external

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## run: run from code.
run:
	go run cmd/distribyted/main.go examples/conf_example.yaml

## build: build binary.
build: go-generate go-build

## test-race: execute all tests with race enabled.
test-race:
	go test -v --race -coverprofile=coverage.out -covermode atomic ./...

## test: execute all tests
test:
	go test -v -coverprofile=coverage.out -covermode atomic ./...

go-build:
	@echo "  >  Building binary on $(BIN_OUTPUT)..."
	go build -o $(BIN_OUTPUT) -tags "release" -ldflags='$(LDFLAGS)' cmd/distribyted/main.go

go-generate:
	@echo "  >  Generating code files..."
	go generate ./...

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
