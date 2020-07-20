@echo off

:: SET DevDir=D:\CoolQ Pro\dev\cn.miaoscraft.sis

echo Setting proxy
SET GOPROXY=https://goproxy.cn

echo Checking go installation...
go version > nul
IF ERRORLEVEL 1 (
	echo Please install go first...
	goto RETURN
)

echo Checking gcc installation...
gcc --version > nul
IF ERRORLEVEL 1 (
	echo Please install gcc first...
	goto RETURN
)

echo Checking cqcfg installation...
cqcfg -v
IF ERRORLEVEL 1 (
	echo Install cqcfg...
	go get github.com/Tnze/CoolQ-Golang-SDK/tools/cqcfg@master
	IF ERRORLEVEL 1 (
		echo Install cqcfg fail
		goto RETURN
	)
)

echo Generating app.json ...
go generate
IF ERRORLEVEL 1 (
	echo Generate app.json fail
	goto RETURN
)
echo.

echo Setting env vars..
SET CGO_LDFLAGS=-Wl,--kill-at
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=386

echo Building app.dll ...
go build -ldflags "-s -w" -buildmode=c-shared -ldflags "-extldflags ""-static""" -o app.dll
IF ERRORLEVEL 1 (pause) ELSE (echo Build success!)

if defined DevDir (
    echo Copy app.dll amd app.json ...
    for %%f in (app.dll,app.json) do move %%f "%DevDir%\%%f" > nul
    IF ERRORLEVEL 1 pause
)

exit /B

:RETURN
if not defined NOPAUSE pause
exit /B
