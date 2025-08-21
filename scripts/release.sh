#!/bin/bash

# 自动发布脚本
# 执行发布前检查，创建标签，并推送到远程仓库

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 显示帮助
show_help() {
    cat << EOF
ENVM 自动发布脚本

用法:
    $0 <版本号> [选项]

参数:
    版本号         要发布的版本号 (例如: v1.0.0, v1.2.3-beta.1)

选项:
    -f, --force    跳过确认提示，强制发布
    -d, --dry-run  只运行检查，不实际创建标签和推送
    -h, --help     显示帮助信息

示例:
    $0 v1.0.0                 # 发布 v1.0.0
    $0 v1.1.0-beta.1 --dry-run  # 预览发布 v1.1.0-beta.1
    $0 v2.0.0 --force         # 强制发布 v2.0.0

发布流程:
    1. 运行发布前检查
    2. 生成变更日志
    3. 确认发布信息
    4. 创建 Git 标签
    5. 推送到远程仓库
    6. GitHub Actions 自动构建和发布
EOF
}

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 验证版本号格式
validate_version() {
    local version="$1"
    
    # 检查版本号格式 (vX.Y.Z 或 vX.Y.Z-suffix)
    if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)*$ ]]; then
        log_error "版本号格式无效: $version"
        log_info "期望格式: vX.Y.Z 或 vX.Y.Z-suffix (例如: v1.0.0, v1.2.3-beta.1)"
        return 1
    fi
    
    # 检查版本号是否已存在
    if git rev-parse "$version" >/dev/null 2>&1; then
        log_error "版本标签 $version 已存在"
        return 1
    fi
    
    return 0
}

# 生成变更日志预览
generate_changelog_preview() {
    local new_version="$1"
    local previous_tag
    
    # 获取前一个标签
    previous_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -z "$previous_tag" ]; then
        previous_tag=$(git rev-list --max-parents=0 HEAD)
        log_warning "这是第一个版本，将显示所有提交"
    fi
    
    echo ""
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${BLUE}        变更日志预览 ($previous_tag -> $new_version)${NC}"
    echo -e "${BLUE}===========================================${NC}"
    
    # 统计信息
    local commit_count
    commit_count=$(git rev-list --count "$previous_tag".."HEAD")
    echo -e "${GREEN}提交数量: $commit_count${NC}"
    
    if [ "$commit_count" -eq 0 ]; then
        log_warning "自上次发布以来没有新的提交"
        return 1
    fi
    
    # 显示提交摘要
    echo ""
    echo "📝 主要变更:"
    git log "$previous_tag"..HEAD --oneline --no-merges | head -10 | sed 's/^/  - /'
    
    if [ "$commit_count" -gt 10 ]; then
        echo "  ... 以及其他 $((commit_count - 10)) 个提交"
    fi
    
    # 显示相关的 Issues 和 PRs
    echo ""
    echo "🔗 相关问题:"
    local issues
    issues=$(git log "$previous_tag"..HEAD --oneline | grep -oE "#[0-9]+" | sort -u | head -5)
    if [ -n "$issues" ]; then
        echo "$issues" | sed 's/^/  - /'
    else
        echo "  - 无直接关联的 Issues"
    fi
    
    echo ""
}

# 确认发布
confirm_release() {
    local version="$1"
    
    echo -e "${YELLOW}===========================================${NC}"
    echo -e "${YELLOW}            确认发布信息${NC}"
    echo -e "${YELLOW}===========================================${NC}"
    echo -e "${GREEN}版本号:${NC} $version"
    echo -e "${GREEN}分支:${NC} $(git rev-parse --abbrev-ref HEAD)"
    echo -e "${GREEN}提交:${NC} $(git rev-parse --short HEAD)"
    echo -e "${GREEN}远程:${NC} $(git remote get-url origin)"
    echo ""
    
    if [ "$FORCE" = "true" ]; then
        log_info "强制模式，跳过确认"
        return 0
    fi
    
    read -p "确认发布 $version? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "用户取消发布"
        return 1
    fi
    
    return 0
}

# 创建发布
create_release() {
    local version="$1"
    
    log_info "创建标签 $version..."
    
    # 创建带注释的标签
    local tag_message="Release $version

$(git log $(git describe --tags --abbrev=0 2>/dev/null || echo "")..HEAD --oneline --no-merges | head -5)

发布时间: $(date -u '+%Y-%m-%d %H:%M:%S UTC')
发布者: $(git config user.name) <$(git config user.email)>
"
    
    git tag -a "$version" -m "$tag_message"
    log_success "标签创建成功"
    
    if [ "$DRY_RUN" = "true" ]; then
        log_info "预演模式，不推送标签"
        log_info "要推送标签，请运行: git push origin $version"
    else
        log_info "推送标签到远程仓库..."
        git push origin "$version"
        log_success "标签推送成功"
        
        echo ""
        log_success "🎉 发布完成！"
        echo ""
        log_info "GitHub Actions 将自动:"
        echo "  - 构建多平台二进制文件"
        echo "  - 生成 Windows 安装程序"
        echo "  - 创建 GitHub Release"
        echo "  - 生成完整的发布说明"
        echo ""
        log_info "您可以在以下位置查看进度:"
        echo "  - Actions: https://github.com/$(git remote get-url origin | sed 's/.*github.com[/:]//;s/.git$//')/actions"
        echo "  - Releases: https://github.com/$(git remote get-url origin | sed 's/.*github.com[/:]//;s/.git$//')/releases"
    fi
}

# 清理函数（发生错误时）
cleanup() {
    local version="$1"
    
    if git rev-parse "$version" >/dev/null 2>&1; then
        log_warning "清理创建的标签 $version"
        git tag -d "$version" 2>/dev/null || true
    fi
}

# 主函数
main() {
    local version="$1"
    
    # 检查必需的工具
    if ! command -v git &> /dev/null; then
        log_error "Git 未安装或不在 PATH 中"
        exit 1
    fi
    
    # 检查是否在 Git 仓库中
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        log_error "当前目录不是 Git 仓库"
        exit 1
    fi
    
    # 验证版本号
    if ! validate_version "$version"; then
        exit 1
    fi
    
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${BLUE}     ENVM 自动发布脚本 - $version${NC}"
    echo -e "${BLUE}===========================================${NC}"
    echo ""
    
    # 运行发布前检查
    log_info "运行发布前检查..."
    if [ -f "./scripts/pre-release-check.sh" ]; then
        if ! bash "./scripts/pre-release-check.sh" "$version"; then
            log_error "发布前检查失败"
            exit 1
        fi
    else
        log_warning "发布前检查脚本不存在，跳过"
    fi
    
    # 生成变更日志预览
    if ! generate_changelog_preview "$version"; then
        log_error "无法生成变更日志"
        exit 1
    fi
    
    # 确认发布
    if ! confirm_release "$version"; then
        log_info "发布已取消"
        exit 0
    fi
    
    # 创建发布
    trap 'cleanup "$version"' ERR
    create_release "$version"
    trap - ERR
}

# 解析命令行参数
FORCE=false
DRY_RUN=false
VERSION=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--force)
            FORCE=true
            shift
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        -*)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
        *)
            if [ -z "$VERSION" ]; then
                VERSION="$1"
            else
                log_error "只能指定一个版本号"
                exit 1
            fi
            shift
            ;;
    esac
done

# 检查是否提供了版本号
if [ -z "$VERSION" ]; then
    log_error "请提供版本号"
    show_help
    exit 1
fi

# 运行主函数
main "$VERSION"