# envm Windows 安装脚本
param(
    [string]$InstallPath = "$env:LOCALAPPDATA\envm",
    [switch]$Force = $false
)

# 颜色输出函数
function Write-ColorOutput {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Message,
        [string]$ForegroundColor = "White"
    )
    Write-Host $Message -ForegroundColor $ForegroundColor
}

# 错误处理
$ErrorActionPreference = "Stop"

try {
    Write-ColorOutput "envm Windows 自动安装脚本" -ForegroundColor Green
    Write-Host ""

    # 检测系统架构
    $arch = if ([Environment]::Is64BitProcess) { "amd64" } else { "386" }
    Write-ColorOutput "检测到系统架构: windows-$arch" -ForegroundColor Green

    # 获取最新版本
    Write-ColorOutput "获取最新版本信息..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/FirewineXie/envm/releases/latest" -UseBasicParsing
        $version = $response.tag_name
        Write-ColorOutput "最新版本: $version" -ForegroundColor Green
    }
    catch {
        Write-ColorOutput "无法获取最新版本信息，使用默认版本" -ForegroundColor Red
        $version = "v1.0.0"
    }

    # 创建安装目录
    if (-not (Test-Path $InstallPath)) {
        Write-ColorOutput "创建安装目录: $InstallPath" -ForegroundColor Yellow
        New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
    }

    # 下载二进制文件
    $downloadUrl = "https://github.com/FirewineXie/envm/releases/download/$version/envm-$version-windows-$arch.zip"
    $zipFile = Join-Path $InstallPath "envm.zip"
    
    Write-ColorOutput "下载 envm $version for windows-$arch..." -ForegroundColor Yellow
    Write-ColorOutput "下载地址: $downloadUrl" -ForegroundColor Gray
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipFile -UseBasicParsing
        Write-ColorOutput "下载完成" -ForegroundColor Green
    }
    catch {
        Write-ColorOutput "下载失败: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }

    # 解压文件
    Write-ColorOutput "解压文件..." -ForegroundColor Yellow
    try {
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        [System.IO.Compression.ZipFile]::ExtractToDirectory($zipFile, $InstallPath)
        Remove-Item $zipFile
        Write-ColorOutput "解压完成" -ForegroundColor Green
    }
    catch {
        Write-ColorOutput "解压失败: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }

    # 设置环境变量
    Write-ColorOutput "设置环境变量..." -ForegroundColor Yellow
    
    $gvmHome = "$env:USERPROFILE\.govm"
    $gvmSymlink = "$env:USERPROFILE\.govm\go"
    
    # 创建 .govm 目录
    if (-not (Test-Path $gvmHome)) {
        New-Item -ItemType Directory -Path $gvmHome -Force | Out-Null
    }

    # 设置用户环境变量
    try {
        [Environment]::SetEnvironmentVariable("GOVM_HOME", $gvmHome, "User")
        [Environment]::SetEnvironmentVariable("GOVM_SYMLINK", $gvmSymlink, "User")
        Write-ColorOutput "环境变量设置完成" -ForegroundColor Green
    }
    catch {
        Write-ColorOutput "设置环境变量失败: $($_.Exception.Message)" -ForegroundColor Red
    }

    # 添加到 PATH
    try {
        $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        
        # 检查是否已在 PATH 中
        if ($userPath -notlike "*$InstallPath*") {
            $newPath = if ($userPath) { "$userPath;$InstallPath" } else { $InstallPath }
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
            Write-ColorOutput "已将 $InstallPath 添加到 PATH" -ForegroundColor Green
        } else {
            Write-ColorOutput "$InstallPath 已在 PATH 中" -ForegroundColor Yellow
        }

        # 同时添加 GOVM_SYMLINK\bin 到 PATH
        if ($userPath -notlike "*$gvmSymlink\bin*") {
            $newPath = [Environment]::GetEnvironmentVariable("PATH", "User")
            $newPath += ";$gvmSymlink\bin"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
            Write-ColorOutput "已将 $gvmSymlink\bin 添加到 PATH" -ForegroundColor Green
        }
    }
    catch {
        Write-ColorOutput "添加到 PATH 失败: $($_.Exception.Message)" -ForegroundColor Red
        Write-ColorOutput "请手动将 $InstallPath 添加到系统 PATH 环境变量" -ForegroundColor Yellow
    }

    # 验证安装
    Write-ColorOutput "验证安装..." -ForegroundColor Yellow
    
    # 刷新当前会话的环境变量
    $env:PATH = [Environment]::GetEnvironmentVariable("PATH", "User") + ";" + [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $env:GOVM_HOME = $gvmHome
    $env:GOVM_SYMLINK = $gvmSymlink
    
    $envmPath = Join-Path $InstallPath "envm-$version-windows-$arch.exe"
    if (Test-Path $envmPath) {
        # 重命名为简短名称
        $finalPath = Join-Path $InstallPath "envm.exe"
        Move-Item $envmPath $finalPath -Force
        
        Write-ColorOutput "envm 安装成功！" -ForegroundColor Green
        Write-ColorOutput "安装路径: $finalPath" -ForegroundColor Green
        
        # 尝试显示版本
        try {
            $versionOutput = & $finalPath --version 2>&1
            Write-ColorOutput "版本: $versionOutput" -ForegroundColor Green
        }
        catch {
            Write-ColorOutput "无法获取版本信息，但安装成功" -ForegroundColor Yellow
        }
    } else {
        Write-ColorOutput "安装验证失败，找不到可执行文件" -ForegroundColor Red
        exit 1
    }

    Write-Host ""
    Write-ColorOutput "安装完成！" -ForegroundColor Green
    Write-Host ""
    Write-ColorOutput "使用方法:" -ForegroundColor Cyan
    Write-Host "  envm list           # 列出可用版本"
    Write-Host "  envm install 1.21   # 安装指定版本"
    Write-Host "  envm use 1.21       # 切换到指定版本"
    Write-Host ""
    Write-ColorOutput "注意事项:" -ForegroundColor Yellow
    Write-Host "1. 请重新启动 PowerShell 或命令提示符以刷新环境变量"
    Write-Host "2. 如果命令无法识别，请检查 PATH 环境变量是否正确设置"
    Write-Host "3. GOVM_HOME: $gvmHome"
    Write-Host "4. GOVM_SYMLINK: $gvmSymlink"

} catch {
    Write-ColorOutput "安装失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}