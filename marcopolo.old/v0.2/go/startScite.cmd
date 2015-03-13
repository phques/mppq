@echo off

REM won't work if dir path has spaces
set GOPATH=%~d0%~p0

start "scite" "C:\Program Files\Util\wscite\SciTE.exe" 

