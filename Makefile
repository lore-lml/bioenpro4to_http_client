ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	cd lib/channel_manager && cargo build
	cp lib/channel_manager/target/debug/libc_channel_manager_lib.so lib/
	go build -ldflags="-r $(ROOT_DIR)lib" main.go

run: build
	./main
