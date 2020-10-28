#-include .env

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
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

## test: execute all tests.
test:
	go test -v --race -coverprofile=coverage.out ./...

compile: go-generate go-compile

go-build:
	@echo "  >  Building binary..."
	go build -o bin/distribyted -tags "release" cmd/distribyted/main.go

go-generate:
	@echo "  >  Generating code files..."
	go generate ./...

go-compile:
	GOPATH=~/go
	ORGPATH=$(GOPATH)/src/github.com/distribyted
	REPOPATH=$(ORGPATH)/distribyted

	@echo "  >  Compiling for several platforms..."
	go install src.techknowlogick.com/xgo
	docker build ./build_tools/ -t distribyted/xgo-cgofuse
	mkdir -p $(ORGPATH)
	ln -sfrnT . $(REPOPATH)

	@echo "  >  Compiling for windows..."
	GOPATH=$(GOPATH) xgo -out bin/distribyted-$(VERSION) -image=distribyted/xgo-cgofuse -ldflags='-extldflags "-static" $(LDFLAGS)' -tags="release" -targets=windows/amd64 $(REPOPATH)/cmd/distribyted/
	@echo "  >  Compiling for linux..."
	GOPATH=$(GOPATH) xgo -out bin/distribyted-$(VERSION) -image=distribyted/xgo-cgofuse -ldflags='$(LDFLAGS)' -tags="release" -targets=linux/arm-7,linux/amd64 $(REPOPATH)/cmd/distribyted/
#	@echo "  >  Compiling for darwin..."
#	GOPATH=$(GOPATH) xgo -out bin/distribyted-$(VERSION) -image=distribyted/xgo-cgofuse -ldflags='$(LDFLAGS)' -tags="release" -targets=darwin/* $(REPOPATH)/cmd/distribyted/

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
