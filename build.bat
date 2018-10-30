@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
echo build for %GOOS%^<%GOARCH%^>
go build -o billing.exe .

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
echo build for %GOOS%^<%GOARCH%^>
go build -o billing .
echo run test ...
billing.exe