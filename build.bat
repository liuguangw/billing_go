@echo off
SET CGO_ENABLED=0
SET GOARCH=386

SET GOOS=windows
echo build for %GOOS%^<%GOARCH%^>
go build -o billing.exe .

SET GOOS=linux
echo build for %GOOS%^<%GOARCH%^>
go build -o billing .

echo build complete
pause