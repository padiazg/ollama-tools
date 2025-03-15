# Makefile

.PHONY: build-tools


build-tools: pkg=github.com/padiazg/ollama-tools/models/version
build-tools: ldflags = -X $(pkg).version=$(shell git describe --tags --always --dirty) 
build-tools: ldflags += -X $(pkg).commit=$(shell git rev-parse HEAD)
build-tools: ldflags += -X $(pkg).buildDate=$(shell date -Iseconds)
build-tools: ldflags += -X $(pkg).buildBy=make

build-tools:
	@echo "Building ollama-tools..."
	@go build -ldflags "$(ldflags)"