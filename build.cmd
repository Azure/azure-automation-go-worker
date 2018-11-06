@echo off
set BINDIR=bin
set BIN_WORKER=worker.exe
set BIN_SANDBOX=sandbox.exe


REM clean
rmdir /s /q %BINDIR%

REM build
set GOOS=windows
set GOARCH=amd64
go build -v -o %BINDIR%\%BIN_WORKER% .\main\worker
go build -v -o %BINDIR%\%BIN_SANDBOX% .\main\sandbox