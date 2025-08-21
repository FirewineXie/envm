# 🚀 ENVM 发布指南

本文档描述了 envm 项目的完整发布流程和自动化系统。

## 📋 发布系统概览

### 🤖 自动化功能

1. **自动构建** - 多平台二进制文件 (Windows/Linux/macOS)
2. **Windows 安装程序** - NSIS 制作的专业安装程序
3. **自动发布说明** - 基于 git 提交自动生成
4. **问题关联** - 自动识别和链接相关 Issues/PRs
5. **完整文档** - 包含安装说明、使用指南等

### 📁 相关文件

```
envm/
├── .github/
│   ├── workflows/
│   │   ├── release.yml                    # 主发布流程
│   │   └── generate-release-notes.yml     # 发布说明生成
│   ├── ISSUE_TEMPLATE/                    # Issue 模板
│   │   ├── bug_report.yml
│   │   ├── feature_request.yml
│   │   ├── question.yml
│   │   └── config.yml
│   └── PULL_REQUEST_TEMPLATE.md           # PR 模板
├── scripts/
│   ├── generate-changelog.sh              # 变更日志生成
│   ├── pre-release-check.sh               # 发布前检查
│   └── release.sh                         # 自动发布脚本
├── envm-installer.nsi                     # NSIS 安装程序脚本
├── install.sh                             # Linux/macOS 安装脚本
├── install.ps1                            # Windows PowerShell 安装脚本
└── build-installer.bat                    # 本地构建安装程序
```

## 🛠 发布流程

### 方式1：自动化脚本（推荐）

```bash
# 1. 运行发布前检查
./scripts/pre-release-check.sh v1.0.0

# 2. 自动发布（包含所有检查和确认）
./scripts/release.sh v1.0.0

# 3. 预览发布（不实际创建标签）
./scripts/release.sh v1.0.0 --dry-run
```

### 方式2：手动发布

```bash
# 1. 确保代码已提交并推送
git status
git push

# 2. 运行发布前检查
./scripts/pre-release-check.sh v1.0.0

# 3. 创建并推送标签
git tag v1.0.0
git push origin v1.0.0
```

### 🎯 GitHub Actions 自动处理

推送标签后，GitHub Actions 将自动：

1. **构建多平台二进制文件**
   - `envm-v1.0.0-windows-amd64.zip`
   - `envm-v1.0.0-linux-amd64.tar.gz`
   - `envm-v1.0.0-linux-arm64.tar.gz`
   - `envm-v1.0.0-darwin-amd64.tar.gz`
   - `envm-v1.0.0-darwin-arm64.tar.gz`

2. **创建 Windows 安装程序**
   - `envm-installer-v1.0.0.exe`

3. **生成完整的发布说明**
   - 自动分类的变更日志
   - 安装说明和使用指南
   - 相关 Issues 和 PRs
   - 贡献者统计
   - 版本统计信息

4. **创建 GitHub Release**
   - 自动上传所有文件
   - 发布完整的说明文档
   - 标记为最新版本

## 📝 提交规范

为了更好地生成发布说明，建议使用以下提交格式：

### 🏷 提交类型

| 类型 | 描述 | 示例 |
|------|------|------|
| `feat:` | 新功能 | `feat: 添加 Java 版本管理支持` |
| `fix:` | 错误修复 | `fix: 修复 Windows 路径解析问题` |
| `docs:` | 文档更新 | `docs: 更新安装说明` |
| `refactor:` | 代码重构 | `refactor: 优化版本检测逻辑` |
| `perf:` | 性能优化 | `perf: 提升下载速度` |
| `test:` | 测试相关 | `test: 添加单元测试` |
| `build:` | 构建相关 | `build: 更新依赖项` |

### 📌 关联 Issues

在提交消息中使用以下格式关联问题：

```bash
# 关闭问题
git commit -m "fix: 修复安装路径问题 (fixes #123)"

# 引用问题
git commit -m "feat: 添加新功能 (ref #456)"
```

## 🎨 发布说明特性

### 📊 自动生成的内容

1. **分类变更日志**
   - 🚀 新功能
   - 🐛 错误修复
   - ✨ 改进优化
   - 📚 文档更新

2. **详细的安装指南**
   - Windows EXE 安装程序
   - 跨平台脚本安装
   - 手动安装说明

3. **问题和 PR 关联**
   - 自动检测提交中的 #数字
   - 链接到相关问题和拉取请求

4. **统计信息**
   - 提交数量
   - 文件变更统计
   - 贡献者列表

### 🎯 发布说明示例

```markdown
# 🎉 envm v1.0.0 发布

## 📦 安装方法

### Windows 用户
#### 🎯 方式1：EXE安装程序（推荐）
1. 下载 `envm-installer-v1.0.0.exe`
2. 双击运行安装程序
...

## 📝 本次更新内容

### 🚀 新功能
- feat: 添加自定义安装路径支持
- feat: 新增桌面快捷方式选项

### 🐛 错误修复
- fix: 修复 Windows 环境变量设置问题
- fix: 解决路径包含空格的问题
...
```

## 🔧 版本号规范

使用 [语义化版本](https://semver.org/lang/zh-CN/)：

- `v1.0.0` - 主要版本（破坏性变更）
- `v1.1.0` - 次要版本（新功能，向后兼容）
- `v1.1.1` - 修订版本（错误修复，向后兼容）
- `v1.1.0-beta.1` - 预发布版本

## 🚨 发布前检查清单

### ✅ 代码质量
- [ ] 所有测试通过
- [ ] 代码格式正确 (`go fmt`)
- [ ] 静态分析通过 (`go vet`)
- [ ] 无明显的 lint 问题

### ✅ 构建验证
- [ ] 本地构建成功
- [ ] 跨平台编译成功
- [ ] 依赖项整洁 (`go mod tidy`)

### ✅ 文档更新
- [ ] README.md 更新
- [ ] 版本号更新
- [ ] 变更说明准备

### ✅ Git 状态
- [ ] 工作目录干净
- [ ] 在正确的分支
- [ ] 版本标签不存在

## 🛠 本地开发和测试

### 构建 Windows 安装程序

```cmd
# 使用批处理脚本
build-installer.bat 1.0.0

# 或手动构建
go build -o envm-1.0.0-windows-amd64.exe .
makensis envm-installer.nsi
```

### 测试发布流程

```bash
# 预览发布（不创建标签）
./scripts/release.sh v1.0.0-test --dry-run

# 本地生成变更日志
./scripts/generate-changelog.sh v1.0.0 v0.9.0
```

## 🐛 故障排除

### NSIS 安装失败

如果 GitHub Actions 中 NSIS 安装失败：

1. 检查 Chocolatey 是否可用
2. 尝试使用预编译的 Docker 镜像
3. 使用 GitHub 缓存优化安装时间

### 发布失败

1. 检查 GitHub Token 权限
2. 验证标签格式是否正确
3. 确认所有必需的文件都存在

## 📞 获取帮助

- 📖 **文档问题**: 更新 README.md 或本文档
- 🐛 **发布问题**: 提交 Issue 并标记 `release` 标签
- 💬 **讨论**: 使用 GitHub Discussions

---

🤖 *此文档随项目持续更新*