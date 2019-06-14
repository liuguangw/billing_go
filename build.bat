@echo off
SET CGO_ENABLED=0
SET GOARCH=386

SET GOOS=windows
echo build for %GOOS%^<%GOARCH%^>
go build billing

SET GOOS=linux
echo build for %GOOS%^<%GOARCH%^>
go build billing

echo build complete
pause