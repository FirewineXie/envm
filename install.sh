#!/bin/bash
# envm 安装脚本 for Linux/macOS

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检测操作系统和架构
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo -e "${RED}不支持的架构: $ARCH${NC}"; exit 1 ;;
    esac
    
    case $OS in
        linux) OS="linux" ;;
        darwin) OS="darwin" ;;
        *) echo -e "${RED}不支持的操作系统: $OS${NC}"; exit 1 ;;
    esac
}

# 获取最新版本
get_latest_version() {
    echo -e "${YELLOW}获取最新版本信息...${NC}"
    VERSION=$(curl -s https://api.github.com/repos/FirewineXie/envm/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo -e "${RED}无法获取最新版本信息${NC}"
        exit 1
    fi
    echo -e "${GREEN}最新版本: $VERSION${NC}"
}

# 下载和安装
install_envm() {
    local download_url="https://github.com/FirewineXie/envm/releases/download/$VERSION/envm-$VERSION-$OS-$ARCH.tar.gz"
    local tmp_dir=$(mktemp -d)
    local tmp_file="$tmp_dir/envm.tar.gz"
    
    echo -e "${YELLOW}下载 envm $VERSION for $OS-$ARCH...${NC}"
    if ! curl -L -o "$tmp_file" "$download_url"; then
        echo -e "${RED}下载失败${NC}"
        rm -rf "$tmp_dir"
        exit 1
    fi
    
    echo -e "${YELLOW}解压文件...${NC}"
    tar -xzf "$tmp_file" -C "$tmp_dir"
    
    # 确保 /usr/local/bin 存在
    if [ ! -d "/usr/local/bin" ]; then
        echo -e "${YELLOW}创建 /usr/local/bin 目录...${NC}"
        sudo mkdir -p /usr/local/bin
    fi
    
    # 安装二进制文件
    echo -e "${YELLOW}安装 envm 到 /usr/local/bin...${NC}"
    sudo mv "$tmp_dir/envm-$VERSION-$OS-$ARCH" /usr/local/bin/envm
    sudo chmod +x /usr/local/bin/envm
    
    # 清理临时文件
    rm -rf "$tmp_dir"
    
    echo -e "${GREEN}envm 安装成功！${NC}"
}

# 设置环境变量
setup_environment() {
    echo -e "${YELLOW}设置环境变量...${NC}"
    
    local govm_home="$HOME/.govm"
    local govm_symlink="$HOME/.govm/go"
    
    # 创建 .govm 目录
    mkdir -p "$govm_home"
    
    # 检测 shell 配置文件
    local shell_rc=""
    if [ -n "$ZSH_VERSION" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        shell_rc="$HOME/.bashrc"
    else
        shell_rc="$HOME/.profile"
    fi
    
    # 检查是否已经设置了环境变量
    if ! grep -q "GOVM_HOME" "$shell_rc" 2>/dev/null; then
        echo -e "${YELLOW}添加环境变量到 $shell_rc${NC}"
        {
            echo ""
            echo "# envm environment variables"
            echo "export GOVM_HOME=\"$govm_home\""
            echo "export GOVM_SYMLINK=\"$govm_symlink\""
            echo "export PATH=\"\$GOVM_SYMLINK/bin:\$PATH\""
        } >> "$shell_rc"
        
        echo -e "${GREEN}环境变量已添加到 $shell_rc${NC}"
        echo -e "${YELLOW}请运行以下命令以重新加载配置：${NC}"
        echo -e "${GREEN}source $shell_rc${NC}"
    else
        echo -e "${GREEN}环境变量已存在${NC}"
    fi
    
    # 设置当前会话的环境变量
    export GOVM_HOME="$govm_home"
    export GOVM_SYMLINK="$govm_symlink"
    export PATH="$govm_symlink/bin:$PATH"
}

# 验证安装
verify_installation() {
    echo -e "${YELLOW}验证安装...${NC}"
    if command -v envm >/dev/null 2>&1; then
        echo -e "${GREEN}envm 安装成功！${NC}"
        echo -e "${GREEN}版本: $(envm --version)${NC}"
        echo ""
        echo -e "${GREEN}使用方法:${NC}"
        echo -e "  envm list           # 列出可用版本"
        echo -e "  envm install 1.21   # 安装指定版本"
        echo -e "  envm use 1.21       # 切换到指定版本"
        echo ""
        echo -e "${YELLOW}注意: 请确保重新加载你的 shell 配置文件或重新开启终端${NC}"
    else
        echo -e "${RED}安装验证失败，请检查 PATH 设置${NC}"
        exit 1
    fi
}

# 主函数
main() {
    echo -e "${GREEN}envm 自动安装脚本${NC}"
    echo ""
    
    detect_platform
    echo -e "${GREEN}检测到平台: $OS-$ARCH${NC}"
    
    get_latest_version
    install_envm
    setup_environment
    verify_installation
    
    echo -e "${GREEN}安装完成！${NC}"
}

main "$@"