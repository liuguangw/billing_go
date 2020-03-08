@echo off
SET CGO_ENABLED=0
SET GOARCH=386
SET GO111MODULE=on
SET GOPROXY=https://goproxy.cn
SET PROJECT_NAME=billing


FOR %%a IN (linux,windows) DO (
  call :buildCommand %%a
)
echo build complete
pause
exit

REM --------------
REM ------编译命令
REM --------------
:buildCommand
SET GOOS=%1
SET TARGET_FILE=%PROJECT_NAME%
echo build for ^<%GOOS%/%GOARCH%^>
if "%GOOS%" == "windows" (
    SET TARGET_FILE=%TARGET_FILE%.exe
) else if not "%GOOS%" == "linux" (
    SET TARGET_FILE=%TARGET_FILE%_%GOOS%
)
go build -ldflags "-s -w" -a -o %TARGET_FILE% .
if %ERRORLEVEL% NEQ 0 (
	pause
	exit
)
goto :eof