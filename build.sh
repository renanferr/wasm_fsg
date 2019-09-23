#!/bin/sh
GOOS=js GOARCH=wasm \
    go build -o public/main.wasm main.go