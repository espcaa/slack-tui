#!/bin/bash
export GOOS=darwin
export GOARCH=amd64
cd ..
echo " Building ..."
go build -o build/macos-amd64/slacktui
echo " Build done: build/macos-amd64/slacktui"
