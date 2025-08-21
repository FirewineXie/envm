# 写一个类似于nvm 的go 开源工具

> 基于GoLand 的切换环境的功能

> 在go module 的基础上进行切换，没有1.13 下面的不进行考虑

## 安装过程

### 自动安装（推荐）

#### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/FirewineXie/envm/main/install.sh | bash
```

#### Windows

**方式1：PowerShell脚本安装（推荐）**
```powershell
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/FirewineXie/envm/main/install.ps1" -UseBasicParsing).Content
```

**方式2：使用.exe安装程序**
1. 从 [Releases页面](https://github.com/FirewineXie/envm/releases) 下载最新的 `envm-installer-x.x.x.exe`（如果可用）
2. 双击运行安装程序（无需管理员权限）
3. 选择安装路径和组件，完成安装

**方式3：便携版**
1. 从 [Releases页面](https://github.com/FirewineXie/envm/releases) 下载 `envm-x.x.x-windows-amd64.zip`
2. 解压到任意目录（建议：`C:\Users\%USERNAME%\envm`）
3. 手动添加到系统PATH环境变量
4. 设置环境变量：
   - `GOVM_HOME=C:\Users\%USERNAME%\.govm`
   - `GOVM_SYMLINK=C:\Users\%USERNAME%\.govm\go`

### 手动安装

1. 从 [Releases页面](https://github.com/FirewineXie/envm/releases) 下载对应平台的二进制文件
2. 解压并将可执行文件放到系统PATH中
3. 设置环境变量：
   - `GOVM_HOME` 例如: `C:\Users\username\.govm` (Windows) 或 `$HOME/.govm` (Linux/macOS)
   - `GOVM_SYMLINK` 例如: `C:\Users\username\.govm\go` (Windows) 或 `$HOME/.govm/go` (Linux/macOS)
4. 将 `GOVM_SYMLINK/bin` 添加到系统PATH
5. 运行 `envm --version` 验证安装

## 尾注

感谢 `gvm`,`nvm` 提供的灵感和代码的实现



## 使用方法

安装完成后，可以使用以下命令：

```bash
# 查看版本
envm --version

# 列出所有可用的Go版本
envm list

# 安装指定版本的Go
envm install 1.21.0

# 切换到指定版本
envm use 1.21.0

# 查看当前使用的版本
envm current
```

## 卸载

### Windows (.exe安装程序)
1. 在"控制面板" -> "程序和功能"中找到"envm"
2. 点击"卸载"按照提示完成卸载
3. 或者在开始菜单中找到"envm" -> "卸载envm"

### Linux/macOS
```bash
sudo rm /usr/local/bin/envm
# 手动清理环境变量配置
```

## 本地编译
go build -o envm