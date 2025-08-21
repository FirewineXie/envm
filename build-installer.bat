@echo off
REM Windows 本地构建安装程序脚本

echo 构建 envm Windows 安装程序...

REM 检查是否安装了 Go
go version >nul 2>&1
if errorlevel 1 (
    echo 错误：未检测到 Go 环境，请先安装 Go
    pause
    exit /b 1
)

REM 检查是否安装了 NSIS
where makensis >nul 2>&1
if errorlevel 1 (
    echo 错误：未检测到 NSIS，请先安装 NSIS
    echo.
    echo 安装方法：
    echo 1. 使用 winget: winget install NSIS.NSIS
    echo 2. 使用 choco: choco install nsis
    echo 3. 手动下载: https://nsis.sourceforge.io/Download
    echo.
    pause
    exit /b 1
)

REM 设置版本号（可以通过参数传入）
set VERSION=%1
if "%VERSION%"=="" set VERSION=1.0.0

echo 版本号：%VERSION%

REM 构建 Windows 二进制文件
echo 构建二进制文件...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-s -w -X main.version=%VERSION%" -o "envm-%VERSION%-windows-amd64.exe" .

if errorlevel 1 (
    echo 构建失败！
    pause
    exit /b 1
)

echo 构建成功：envm-%VERSION%-windows-amd64.exe

REM 更新 NSIS 脚本中的版本号
echo 更新 NSIS 脚本版本号...
powershell -Command "(Get-Content envm-installer.nsi) -replace '!define PRODUCT_VERSION \"1.0.0\"', '!define PRODUCT_VERSION \"%VERSION%\"' -replace 'OutFile \"envm-installer-1.0.0.exe\"', 'OutFile \"envm-installer-%VERSION%.exe\"' | Set-Content envm-installer.nsi"

REM 构建安装程序
echo 构建安装程序...
makensis envm-installer.nsi

if errorlevel 1 (
    echo 安装程序构建失败！
    pause
    exit /b 1
)

echo.
echo ==========================================
echo 构建完成！
echo 二进制文件：envm-%VERSION%-windows-amd64.exe
echo 安装程序：  envm-installer-%VERSION%.exe
echo ==========================================
echo.

REM 询问是否清理临时文件
set /p cleanup="是否删除临时二进制文件？(y/n): "
if /i "%cleanup%"=="y" (
    del "envm-%VERSION%-windows-amd64.exe"
    echo 临时文件已清理
)

pause