#!/bin/sh

mkdir "releases"

# 可选参数-ldflags 是编译选项：
#   -s -w 去掉调试信息，可以减小构建后文件体积。

# 【darwin/amd64】
echo "start build darwin/amd64 ..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o ./releases/goshowdoc-darwin-amd64

# 【linux/amd64】
echo "start build linux/amd64 ..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o ./releases/goshowdoc-linux-amd64

# 【windows/amd64】
echo "start build windows/amd64 ..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o ./releases/goshowdoc-windows-amd64.exe

echo "Congratulations,all build success!!!"