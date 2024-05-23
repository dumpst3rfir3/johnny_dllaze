#!/usr/bin/env bash

# Setup ENV
export GOOS=windows

# Version Info
go generate

# Compile EXE for windows using mingw
go build  --buildvcs=false --ldflags="-s -w -X 'main.dll=${1:-updater.dll}'" \
    -o "goader.exe"  --trimpath .

