@ECHO OFF
RD /S /Q release > NUL 2>&1
MKDIR release > NUL 2>&1

ECHO Building Linux x64
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-w -s" -o release\shadowsocks-legendsock

ECHO Building Windows x64
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-w -s" -o release\shadowsocks-legendsock.exe

ECHO UPX Compressing
upx release\*

PAUSE