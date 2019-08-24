@echo off
SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386
SET GOPROXY=https://goproxy.cn

::检查是否安装了编译器
where gcc > nul
if errorlevel 1 echo 找不到gcc
where go > nul
if errorlevel 1 echo 找不到go

::下载SDK  
go get github.com/Tnze/CoolQ-Golang-SDK > nul
::安装cqcfg  
go install github.com/Tnze/CoolQ-Golang-SDK/tools/cqcfg
