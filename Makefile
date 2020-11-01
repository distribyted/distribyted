#-include .env

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

go-cross-compile: GOPATH=~/go
go-cross-compile: ORGPATH=$(GOPATH)/src/github.com/distribyted
go-cross-compile: REPOPATH=$(ORGPATH)/distribyted

# Use linker flags to provide version/build settings
LDFLAGS=-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -linkmode external

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## run: run from code.
run:
	go run cmd/distribyted/main.go examples/conf_example.yaml

## build: build binary.
build: go-generate go-build

## test: execute all tests.
test:
	go test -v --race -coverprofile=coverage.out ./...

## cross-compile: compile for other platforms using xgo.
cross-compile: go-generate go-cross-compile

go-build:
	@echo "  >  Building binary..."
	go build -o bin/distribyted-$(VERSION)-`go env GOOS`-`go env GOARCH``go env GOEXE` -tags "release" -ldflags='$(LDFLAGS)' cmd/distribyted/main.go

go-generate:
	@echo "  >  Generating code files..."
	go generate ./...

go-cross-compile:
	@echo "  >  Compiling for several platforms..."
	GO111MODULE=off go get -u src.techknowlogick.com/xgo
	docker build ./build_tools/ -t distribyted/xgo-cgofuse
	mkdir -p $(ORGPATH)
	ln -sfrnT . $(REPOPATH)

	GOPATH=$(GOPATH) xgo -out bin/distribyted-$(VERSION) -image=distribyted/xgo-cgofuse -ldflags='$(LDFLAGS)' -tags="release" -targets=linux/arm-7 $(REPOPATH)/cmd/distribyted/

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
