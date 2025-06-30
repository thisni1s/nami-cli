#!/bin/sh

GOOS=linux   GOARCH=amd64 go build -o build/nami-cli-linux-amd64
GOOS=windows GOARCH=amd64 go build -o build/nami-cli-windows-amd64.exe
GOOS=darwin  GOARCH=amd64 go build -o build/nami-cli-darwin-amd64
GOOS=darwin  GOARCH=arm64 go build -o build/nami-cli-darwin-arm64
