#!/bin/bash
export GOOS=linux
export GOARCH=amd64
cd ..
echo " Building ..."
go build -o build/linux-amd64/slacktui
echo " Build done: build/linux-amd64/slacktui"
