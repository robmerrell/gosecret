GO ?= go
GXC ?= gxc

all: build

build:
	$(GO) build

build-linux:
	$(GXC) build linux/amd64
