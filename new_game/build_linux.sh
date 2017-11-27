#!/usr/bin/env bash
GOOS=linux
GOARCH=amd64

cd flag_adder && $GOROOT/bin/go build && cd ..
cd flag_handler && $GOROOT/bin/go build && cd ..
cd round_handler && $GOROOT/bin/go build && cd ..
cd router && $GOROOT/bin/go build && cd ..
