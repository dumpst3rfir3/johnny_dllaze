#!/usr/bin/env bash

# Setup ENV
export CC=x86_64-w64-mingw32-gcc
export CGO_ENABLED=1
export GOOS=windows

# Version Info
go generate

# Compile DLL for windows using mingw
go build --buildmode=c-shared --buildvcs=false --ldflags="-s -w" \
    -o "${1:-updater.dll}" --tags=dll --trimpath .
