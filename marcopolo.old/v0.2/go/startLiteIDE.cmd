@echo off

REM won't work if dir path has spaces
set GOPATH=%~d0%~p0

start "liteide" C:\Dev\liteide\bin\liteide.exe %GOPATH%

