#!/bin/bash

ldflags="-s -w"

go env -w GOPROXY=https://goproxy.cn,direct

# Plugin
filepath="plugin"
for item in "$filepath"/*.go; do
  buf=${item%.go}
  name=${buf##*/}
  echo $name
  CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags "$ldflags" -o plugin/score-"$name" plugin/"$name".go
  upx plugin/score-"$name"
done

# Test
target="plugin-score-test"
CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags "$ldflags" -o $target main.go
upx $target
