:: 关闭控制台回显  
@echo off

:: 酷Q的dev文件夹路径（改成你自己的）
SET DevDir=D:\酷Q Pro\dev\me.cqp.tnze.demo
if not exist "%DevDir%" mkdir "%DevDir%"

:: 设置环境变量  
SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386
SET GOPROXY=https://goproxy.cn

:: 生成app.json  
go generate

:: 编译app.dll  
go build -buildmode=c-shared -o app.dll

:: 把app.dll和app.json复制到酷Q的dev文件夹
for %%f in (app.dll,app.json) do move %%f "%DevDir%\%%f" > nul