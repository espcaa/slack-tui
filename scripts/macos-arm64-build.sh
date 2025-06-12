#!/bin/bash
export GOOS=darwin
export GOARCH=arm64
cd ..
echo " Building ..."
go build -o build/macos-arm64/slacktui
echo " Build done: build/macos-arm64/slacktui"
