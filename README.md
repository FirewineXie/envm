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

**方式1：使用.exe安装程序（推荐）**
1. 从 [Releases页面](https://github.com/FirewineXie/envm/releases) 下载最新的 `envm-installer-x.x.x.exe`
2. 双击运行安装程序（无需管理员权限）
3. 选择安装路径（默认：`%LOCALAPPDATA%\envm`）
4. 选择要安装的组件：
   - 核心程序（必需）
   - 开始菜单快捷方式
   - 桌面快捷方式
5. 完成安装（自动设置环境变量和PATH）

**安装程序特性：**
- 🎯 **自定义安装路径** - 可选择任意安装位置
- 🛡️ **路径验证** - 自动检查路径长度和写入权限
- 🔄 **智能升级** - 检测已安装版本，支持覆盖安装
- ⚙️ **组件选择** - 可选安装开始菜单和桌面快捷方式
- 🗑️ **完整卸载** - 通过控制面板完全卸载

**方式2：PowerShell脚本安装**
```powershell
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/FirewineXie/envm/main/install.ps1" -UseBasicParsing).Content
```

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