NAME=weekly-report-gen
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'github.com/tetsuyanh/weekly-report-gen/cmd.version=$(VERSION)' \
           -X 'github.com/tetsuyanh/weekly-report-gen/cmd.revision=$(REVISION)'

## Setup
setup:
	go mod download

## Run application
run: ./main.go
	go run $<

## Build bainary named tag version
build: ./main.go
	go build -ldflags "$(LDFLAGS)" -o bin/$(NAME) $<

## Build each environment binaries
cross-build: ./main.go
	GOOS=linux GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/linux/amd64/$(NAME) $<
	GOOS=darwin GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/darwin/amd64/$(NAME) $<
	GOOS=windows GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/windows/amd64/$(NAME) $<

## Show info
info:
	@echo version: ${VERSION}
	@echo revision: ${REVISION}


.PHONY: setup help
