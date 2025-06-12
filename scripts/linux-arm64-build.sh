#!/bin/bash
export GOOS=linux
export GOARCH=arm64
cd ..
echo " Building ..."
go build -o build/linux-arm64/slacktui
echo " Build done: build/linux-arm64/slacktui"
