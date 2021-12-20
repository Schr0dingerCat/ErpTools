#!/bin/bash
# 清理
go clean
rm -rf ./build/*
# build
gofmt -l -w .
GOOS=windows GOARCH=amd64 go build -o erptools.exe -ldflags="-H windowsgui" main.go
# GOOS=windows GOARCH=amd64 go build -o erptools.exe main.go
# GOOS=windows GOARCH=386 go build -o erptools.exe main.go
mkdir -p ./build
mkdir ./build/config
mkdir ./build/dictionary
mkdir ./build/dist
cp ./config/config.json ./build/config
cp ./config/tw2s.json ./build/config
cp ./dictionary/TW* ./build/dictionary
cp ./dictionary/TS* ./build/dictionary
cp -r ./dist ./build/
mv erptools.exe ./build
