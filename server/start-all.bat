@echo off
chcp 65001 >nul
rem 一键启动 90 天打卡完整链路：MySQL + Go 后端 + Cloudflare Tunnel
rem 用途：电脑重启后双击此文件，等待所有服务就绪即可外网访问

echo [1/3] 检查 MySQL...
tasklist /FI "IMAGENAME eq mysqld.exe" | find /I "mysqld.exe" >nul
if errorlevel 1 (
    echo   正在启动 MySQL...
    start "MySQL" /min "D:\mysql-8.0.46-winx64\bin\mysqld.exe" --console
    timeout /t 5 /nobreak >nul
) else (
    echo   MySQL 已在运行
)

echo [2/3] 检查 Go 后端...
tasklist /FI "IMAGENAME eq roadmap-server.exe" | find /I "roadmap-server.exe" >nul
if errorlevel 1 (
    echo   正在启动 roadmap-server...
    cd /d D:\WorkBuddy\personal-homepage\server
    start "roadmap-server" /min roadmap-server.exe
    timeout /t 3 /nobreak >nul
) else (
    echo   roadmap-server 已在运行
)

echo [3/3] 检查 Cloudflare Tunnel...
rem 如果 tunnel 已注册为 Windows 服务（cloudflared service install），它会开机自启，无需手动启动
sc query Cloudflared | find "RUNNING" >nul 2>&1
if not errorlevel 1 (
    echo   Tunnel 服务已在运行
) else (
    echo   提示：Tunnel 服务未运行。若尚未安装服务，请按 README 执行：
    echo     D:\cloudflared\cloudflared.exe service install 你的TUNNEL_TOKEN
)

echo.
echo 完成！后端地址: http://localhost:8080/api/health
pause
