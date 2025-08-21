#!/bin/bash

# 发布前检查脚本
# 确保所有必要的检查都通过，然后才能发布新版本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检查结果计数
PASSED=0
FAILED=0
WARNINGS=0

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED++))
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    ((WARNINGS++))
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED++))
}

# 检查函数
check_git_status() {
    log_info "检查 Git 工作目录状态..."
    
    if [ -n "$(git status --porcelain)" ]; then
        log_error "工作目录不干净，存在未提交的更改"
        git status --short
        return 1
    else
        log_success "工作目录干净"
    fi
}

check_branch() {
    log_info "检查当前分支..."
    
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
        log_warning "当前不在主分支 ($CURRENT_BRANCH)，确保这是有意的"
    else
        log_success "在主分支 ($CURRENT_BRANCH)"
    fi
}

check_go_version() {
    log_info "检查 Go 版本..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装"
        return 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_success "Go 版本: $GO_VERSION"
    
    # 检查 go.mod 中的 Go 版本要求
    if [ -f "go.mod" ]; then
        MOD_GO_VERSION=$(grep "^go " go.mod | awk '{print $2}')
        log_info "go.mod 要求的版本: $MOD_GO_VERSION"
    fi
}

check_build() {
    log_info "检查代码构建..."
    
    if go build -o /tmp/envm-test .; then
        log_success "构建成功"
        rm -f /tmp/envm-test
    else
        log_error "构建失败"
        return 1
    fi
}

check_tests() {
    log_info "运行测试..."
    
    if go test ./...; then
        log_success "所有测试通过"
    else
        log_error "测试失败"
        return 1
    fi
}

check_format() {
    log_info "检查代码格式..."
    
    UNFORMATTED=$(gofmt -l .)
    if [ -n "$UNFORMATTED" ]; then
        log_error "以下文件格式不正确:"
        echo "$UNFORMATTED"
        log_info "运行 'go fmt ./...' 来修复"
        return 1
    else
        log_success "代码格式正确"
    fi
}

check_lint() {
    log_info "检查代码规范..."
    
    if command -v golint &> /dev/null; then
        LINT_ISSUES=$(golint ./...)
        if [ -n "$LINT_ISSUES" ]; then
            log_warning "发现 lint 问题:"
            echo "$LINT_ISSUES"
        else
            log_success "无 lint 问题"
        fi
    else
        log_warning "golint 未安装，跳过 lint 检查"
    fi
}

check_vet() {
    log_info "运行 go vet..."
    
    if go vet ./...; then
        log_success "go vet 检查通过"
    else
        log_error "go vet 发现问题"
        return 1
    fi
}

check_dependencies() {
    log_info "检查依赖项..."
    
    if go mod tidy -diff; then
        log_success "依赖项整洁"
    else
        log_error "依赖项不整洁，运行 'go mod tidy'"
        return 1
    fi
    
    if go mod verify; then
        log_success "依赖项验证通过"
    else
        log_error "依赖项验证失败"
        return 1
    fi
}

check_version_tag() {
    log_info "检查版本标签..."
    
    if [ -n "$1" ]; then
        NEW_TAG="$1"
        if git rev-parse "$NEW_TAG" >/dev/null 2>&1; then
            log_error "标签 $NEW_TAG 已存在"
            return 1
        else
            log_success "标签 $NEW_TAG 不存在，可以使用"
        fi
        
        # 验证标签格式 (vX.Y.Z)
        if [[ "$NEW_TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)*$ ]]; then
            log_success "标签格式正确: $NEW_TAG"
        else
            log_warning "标签格式可能不标准: $NEW_TAG (期望: vX.Y.Z)"
        fi
    else
        log_warning "未提供新版本标签进行检查"
    fi
}

check_changelog() {
    log_info "检查变更日志..."
    
    # 检查是否有未提交的重要文件
    IMPORTANT_FILES=("README.md" "go.mod" "go.sum")
    for file in "${IMPORTANT_FILES[@]}"; do
        if [ -f "$file" ]; then
            if git diff --quiet HEAD~1 "$file" 2>/dev/null; then
                log_info "$file 无变更"
            else
                log_success "$file 有更新"
            fi
        fi
    done
}

check_cross_platform() {
    log_info "检查跨平台编译..."
    
    PLATFORMS=("windows/amd64" "linux/amd64" "darwin/amd64")
    
    for platform in "${PLATFORMS[@]}"; do
        IFS="/" read -r GOOS GOARCH <<< "$platform"
        if env GOOS="$GOOS" GOARCH="$GOARCH" go build -o /tmp/envm-"$GOOS"-"$GOARCH" . 2>/dev/null; then
            log_success "跨平台编译成功: $platform"
            rm -f /tmp/envm-"$GOOS"-"$GOARCH"*
        else
            log_error "跨平台编译失败: $platform"
            return 1
        fi
    done
}

check_security() {
    log_info "安全检查..."
    
    # 检查是否有硬编码的敏感信息
    if grep -r "password\|secret\|token\|key" --include="*.go" . | grep -v "test" | grep -v "example"; then
        log_warning "发现可能的敏感信息，请检查"
    else
        log_success "未发现明显的敏感信息"
    fi
    
    # 检查 gosec (如果安装了)
    if command -v gosec &> /dev/null; then
        if gosec ./... 2>/dev/null; then
            log_success "安全扫描通过"
        else
            log_warning "安全扫描发现问题，请检查"
        fi
    else
        log_warning "gosec 未安装，跳过安全扫描"
    fi
}

# 主函数
main() {
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${BLUE}       ENVM 发布前检查脚本${NC}"
    echo -e "${BLUE}===========================================${NC}"
    echo ""
    
    NEW_VERSION="$1"
    if [ -n "$NEW_VERSION" ]; then
        log_info "准备发布版本: $NEW_VERSION"
    else
        log_info "运行常规检查 (未指定版本)"
    fi
    
    echo ""
    
    # 运行所有检查
    check_git_status
    check_branch
    check_go_version
    check_format
    check_lint
    check_vet
    check_build
    check_tests
    check_dependencies
    check_cross_platform
    check_security
    check_changelog
    
    if [ -n "$NEW_VERSION" ]; then
        check_version_tag "$NEW_VERSION"
    fi
    
    # 输出检查结果
    echo ""
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${BLUE}              检查结果摘要${NC}"
    echo -e "${BLUE}===========================================${NC}"
    echo -e "${GREEN}通过: $PASSED${NC}"
    echo -e "${YELLOW}警告: $WARNINGS${NC}"
    echo -e "${RED}失败: $FAILED${NC}"
    
    echo ""
    
    if [ $FAILED -eq 0 ]; then
        echo -e "${GREEN}🎉 所有检查通过！可以继续发布流程。${NC}"
        echo ""
        if [ -n "$NEW_VERSION" ]; then
            echo -e "${BLUE}建议的发布步骤:${NC}"
            echo "1. git tag $NEW_VERSION"
            echo "2. git push origin $NEW_VERSION"
            echo "3. GitHub Actions 将自动构建和发布"
        fi
        exit 0
    else
        echo -e "${RED}❌ 发现 $FAILED 个问题，请修复后再发布。${NC}"
        exit 1
    fi
}

# 脚本帮助
show_help() {
    cat << EOF
ENVM 发布前检查脚本

用法:
    $0 [版本号]

参数:
    版本号    可选，要发布的版本号 (例如: v1.0.0)

示例:
    $0                    # 运行常规检查
    $0 v1.0.0            # 检查并准备发布 v1.0.0

检查项目:
    - Git 工作目录状态
    - 当前分支
    - Go 版本和环境
    - 代码格式化
    - 代码规范检查 (lint)
    - 静态分析 (vet)
    - 构建测试
    - 单元测试
    - 依赖项管理
    - 跨平台编译
    - 安全检查
    - 版本标签检查
EOF
}

# 检查参数
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# 运行主函数
main "$1"