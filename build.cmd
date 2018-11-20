@echo off
set BINDIR=bin
set BIN_WORKER=worker.exe
set BIN_SANDBOX=sandbox.exe

go test ./...
if %errorlevel% neq 0 echo Unit test failure. && goto :error
echo All unit tests completed

REM clean
rmdir /s /q %BINDIR%

REM build
set GOOS=windows
set GOARCH=amd64
go build -v -o %BINDIR%\%BIN_WORKER% .\main\worker
go build -v -o %BINDIR%\%BIN_SANDBOX% .\main\sandbox
echo Build complete (binaries can be found in ./%BINDIR%)

:; exit 0
exit /b 0

:error
exit /b %errorlevel%