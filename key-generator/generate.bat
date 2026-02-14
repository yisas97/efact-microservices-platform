@echo off
if not exist bin mkdir bin
javac src\GenerateKeys.java -d bin
if errorlevel 1 exit /b 1
java -cp bin GenerateKeys
if errorlevel 1 exit /b 1
pause
