#!/bin/bash

# 生成 changelog 的脚本
# 用法: ./generate-changelog.sh v1.0.0 v0.9.0

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 获取参数
NEW_TAG=${1:-$(git describe --tags --abbrev=0)}
OLD_TAG=${2:-$(git describe --tags --abbrev=0 ${NEW_TAG}^)}

echo -e "${BLUE}生成从 $OLD_TAG 到 $NEW_TAG 的 changelog${NC}"

# 创建临时文件
TEMP_FILE=$(mktemp)

# 生成 changelog 头部
cat > "$TEMP_FILE" << EOF
## 🚀 新功能 (Features)

EOF

# 获取新功能相关的 commits
git log ${OLD_TAG}..${NEW_TAG} --oneline --grep="feat:" --grep="feature:" --grep="新增" --grep="添加" | while read line; do
    echo "- $line" >> "$TEMP_FILE"
done

# 修复相关
echo "" >> "$TEMP_FILE"
echo "## 🐛 错误修复 (Bug Fixes)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

git log ${OLD_TAG}..${NEW_TAG} --oneline --grep="fix:" --grep="修复" --grep="bugfix" | while read line; do
    echo "- $line" >> "$TEMP_FILE"
done

# 改进优化
echo "" >> "$TEMP_FILE"
echo "## ✨ 改进优化 (Improvements)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

git log ${OLD_TAG}..${NEW_TAG} --oneline --grep="improve:" --grep="refactor:" --grep="优化" --grep="改进" --grep="重构" | while read line; do
    echo "- $line" >> "$TEMP_FILE"
done

# 文档更新
echo "" >> "$TEMP_FILE"
echo "## 📚 文档更新 (Documentation)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

git log ${OLD_TAG}..${NEW_TAG} --oneline --grep="docs:" --grep="doc:" --grep="文档" --grep="readme" | while read line; do
    echo "- $line" >> "$TEMP_FILE"
done

# 其他变更
echo "" >> "$TEMP_FILE"
echo "## 🔧 其他变更 (Others)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

git log ${OLD_TAG}..${NEW_TAG} --oneline --invert-grep --grep="feat:" --grep="fix:" --grep="docs:" --grep="improve:" --grep="refactor:" --grep="新增" --grep="修复" --grep="文档" --grep="优化" --grep="改进" | while read line; do
    echo "- $line" >> "$TEMP_FILE"
done

# 获取相关的 Issues 和 PRs (通过 commit 消息中的 #数字)
echo "" >> "$TEMP_FILE"
echo "## 🔗 相关问题和拉取请求 (Related Issues & PRs)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

# 查找 commit 消息中的 issue/PR 引用
ISSUES=$(git log ${OLD_TAG}..${NEW_TAG} --oneline | grep -oE "#[0-9]+" | sort -u)
if [ ! -z "$ISSUES" ]; then
    for issue in $ISSUES; do
        echo "- $issue" >> "$TEMP_FILE"
    done
else
    echo "- 无直接关联的 Issues 或 PRs" >> "$TEMP_FILE"
fi

# 贡献者统计
echo "" >> "$TEMP_FILE"
echo "## 👥 贡献者 (Contributors)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

git log ${OLD_TAG}..${NEW_TAG} --format='%an' | sort | uniq -c | sort -rn | while read count author; do
    echo "- $author ($count commits)" >> "$TEMP_FILE"
done

# 统计信息
echo "" >> "$TEMP_FILE"
echo "## 📊 版本统计 (Statistics)" >> "$TEMP_FILE"
echo "" >> "$TEMP_FILE"

COMMIT_COUNT=$(git rev-list --count ${OLD_TAG}..${NEW_TAG})
FILES_CHANGED=$(git diff --stat ${OLD_TAG}..${NEW_TAG} | tail -1 | grep -oE '[0-9]+ files? changed' || echo "0 files changed")
INSERTIONS=$(git diff --stat ${OLD_TAG}..${NEW_TAG} | tail -1 | grep -oE '[0-9]+ insertions?' || echo "0 insertions")
DELETIONS=$(git diff --stat ${OLD_TAG}..${NEW_TAG} | tail -1 | grep -oE '[0-9]+ deletions?' || echo "0 deletions")

echo "- **提交数量**: $COMMIT_COUNT" >> "$TEMP_FILE"
echo "- **文件变更**: $FILES_CHANGED" >> "$TEMP_FILE"
echo "- **代码行数**: $INSERTIONS, $DELETIONS" >> "$TEMP_FILE"

# 输出结果
echo -e "${GREEN}Changelog 已生成:${NC}"
cat "$TEMP_FILE"

# 保存到文件
cp "$TEMP_FILE" "CHANGELOG-${NEW_TAG}.md"
echo -e "${YELLOW}Changelog 已保存到 CHANGELOG-${NEW_TAG}.md${NC}"

# 清理
rm "$TEMP_FILE"