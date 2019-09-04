:: 关闭控制台回显  
@echo off

:: 生成app.json
go build github.com/Tnze/CoolQ-Golang-SDK/tools/cqcfg
go generate

:: 设置环境变量  
SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386
SET GOPROXY=https://goproxy.cn

:: 编译app.dll  
go build -buildmode=c-shared -o app.dll

:: 如果设置了环境变量，则把app.dll和app.json复制到酷Q的dev文件夹
REM SET DevDir=D:\酷Q Pro\dev\cn.miaoscraft.sis
if defined DevDir
for %%f in (app.dll,app.json) do move %%f "%DevDir%\%%f" > nul