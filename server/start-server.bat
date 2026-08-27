@echo off
rem 启动 90 天打卡后端（Go + MySQL）
rem 前置：MySQL80 服务已启动（net start MySQL80）
cd /d %~dp0
set GOPROXY=https://goproxy.cn,direct
D:\go\bin\go.exe run main.go
pause
