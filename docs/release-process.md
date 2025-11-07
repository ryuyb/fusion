# Release Process

本文档描述了 Fusion 项目的发布流程，包括版本管理、构建和发布步骤。

## 版本规范

Fusion 项目遵循 [语义化版本（Semantic Versioning）](https://semver.org/lang/zh-CN/)：

```
v{MAJOR}.{MINOR}.{PATCH}[-{PRERELEASE}][+{BUILD}]
```

### 版本号说明

- **MAJOR**：主版本号，当做了不兼容的 API 修改时递增
- **MINOR**：次版本号，当做了向下兼容的功能性新增时递增
- **PATCH**：修订号，当做了向下兼容的问题修正时递增
- **PRERELEASE**：预发布版本标识（可选）
- **BUILD**：构建元数据（可选）

### 示例

- `v1.0.0` - 稳定版本
- `v1.2.3-alpha.1` - Alpha 预发布版本
- `v1.2.3-beta.2` - Beta 预发布版本
- `v1.2.3-rc.1` - Release Candidate 版本
- `v2.0.0-draft` - 草稿版本

## 发布类型

### 1. 稳定版本 (Stable Release)

**适用场景：** 经过充分测试，可以在生产环境使用的版本

**标签格式：** `v{MAJOR}.{MINOR}.{PATCH}`

**特点：**
- GitHub Release 为正式发布
- Docker 镜像会被标记为 `latest`（主版本）
- 生成完整的镜像标签：`v1.2.3`、`1.2`、`1`

**示例：**
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### 2. 预发布版本 (Pre-release)

**适用场景：** 功能完整但需要进一步测试的版本

**标签格式：** `v{MAJOR}.{MINOR}.{PATCH}-{prerelease-identifier}.{number}`

**支持的预发布标识符：**
- `alpha` - 内部测试版本，功能不完整
- `beta` - 公开测试版本，功能完整但可能有 bug
- `rc` (Release Candidate) - 发布候选版本，接近最终版本
- `pre` - 通用预发布标识
- `dev` - 开发版本

**特点：**
- GitHub Release 标记为 Pre-release
- 仅生成完整版本号的镜像标签，不覆盖 major/minor 标签
- 发布说明中包含警告标志

**示例：**
```bash
# Alpha 版本 - 早期测试
git tag -a v1.1.0-alpha.1 -m "Alpha release for testing new features"
git push origin v1.1.0-alpha.1

# Beta 版本 - 功能冻结，进入测试
git tag -a v1.1.0-beta.1 -m "Beta release for public testing"
git push origin v1.1.0-beta.1

# RC 版本 - 发布候选
git tag -a v1.1.0-rc.1 -m "Release candidate 1"
git push origin v1.1.0-rc.1
```

### 3. 草稿版本 (Draft Release)

**适用场景：** 需要内部审查或准备发布说明的版本

**标签格式：** `v{MAJOR}.{MINOR}.{PATCH}-draft`

**特点：**
- GitHub Release 标记为 Draft（草稿）
- 不会公开显示，仅项目成员可见
- 可以继续编辑和完善发布说明

**示例：**
```bash
git tag -a v1.2.0-draft -m "Draft release for review"
git push origin v1.2.0-draft
```

## 完整发布流程

### 步骤 1: 准备发布

1. **确保代码准备就绪**
   ```bash
   # 切换到 main 分支
   git checkout main
   git pull origin main

   # 确保所有测试通过
   make lint
   go test ./...
   ```

2. **更新版本相关文件**（如果有）
   - 更新 CHANGELOG.md
   - 更新文档中的版本号
   - 更新示例配置文件

3. **提交所有更改**
   ```bash
   git add .
   git commit -m "chore: prepare for release v1.0.0"
   git push origin main
   ```

### 步骤 2: 创建标签

根据发布类型选择合适的标签格式：

```bash
# 稳定版本
git tag -a v1.0.0 -m "Release version 1.0.0

- Feature: Add user authentication
- Feature: Implement live streaming
- Fix: Resolve memory leak issue
- Docs: Update API documentation"

# 预发布版本
git tag -a v1.1.0-beta.1 -m "Beta 1 for version 1.1.0

- Feature: New video processing engine (testing)
- Fix: Performance improvements"

# 草稿版本
git tag -a v1.2.0-draft -m "Draft release for review"
```

### 步骤 3: 推送标签

```bash
# 推送单个标签
git push origin v1.0.0

# 或推送所有标签
git push origin --tags
```

### 步骤 4: 自动化构建

推送标签后，GitHub Actions 会自动执行以下操作：

1. **构建多平台二进制文件**
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - 创建压缩包 (`.tar.gz`)

2. **构建 Docker 镜像**
   - 支持多架构：linux/amd64, linux/arm64
   - 推送到 GitHub Container Registry

3. **生成 Changelog**
   - 自动从 Git 提交历史生成
   - 包含自上次标签以来的所有提交

4. **创建 GitHub Release**
   - 上传二进制文件
   - 包含 Docker 镜像拉取命令
   - 显示构建信息和变更日志

### 步骤 5: 验证发布

1. **检查 GitHub Actions**
   - 访问 `https://github.com/{owner}/{repo}/actions`
   - 确认 Release workflow 成功完成

2. **验证 GitHub Release**
   - 访问 `https://github.com/{owner}/{repo}/releases`
   - 检查发布说明和附件

3. **验证 Docker 镜像**
   ```bash
   # 拉取镜像
   docker pull ghcr.io/{owner}/fusion:v1.0.0

   # 验证版本
   docker run --rm ghcr.io/{owner}/fusion:v1.0.0 version
   ```

4. **测试二进制文件**
   ```bash
   # 下载并解压
   wget https://github.com/{owner}/fusion/releases/download/v1.0.0/fusion-linux-amd64.tar.gz
   tar -xzf fusion-linux-amd64.tar.gz

   # 验证版本
   ./fusion-linux-amd64 version
   ```

### 步骤 6: 发布后操作

1. **更新文档**
   - 更新官方文档中的安装说明
   - 更新 Docker Compose 示例中的镜像版本

2. **通知用户**
   - 发布公告（如果是重要版本）
   - 更新社交媒体
   - 通知相关团队

3. **监控反馈**
   - 监控 GitHub Issues
   - 关注用户反馈
   - 准备热修复版本（如需要）

## Docker 镜像标签策略

### 稳定版本标签

发布 `v1.2.3` 时，会创建以下标签：
- `ghcr.io/{owner}/fusion:v1.2.3` - 完整版本号
- `ghcr.io/{owner}/fusion:1.2` - 次版本号
- `ghcr.io/{owner}/fusion:1` - 主版本号
- `ghcr.io/{owner}/fusion:sha-{commit}` - Git 提交 SHA

### 预发布版本标签

发布 `v1.2.3-beta.1` 时，仅创建：
- `ghcr.io/{owner}/fusion:v1.2.3-beta.1` - 完整版本号
- `ghcr.io/{owner}/fusion:sha-{commit}` - Git 提交 SHA

### 开发版本标签

推送到 `main` 分支时，创建：
- `ghcr.io/{owner}/fusion:latest` - 最新开发版本
- `ghcr.io/{owner}/fusion:main` - 主分支
- `ghcr.io/{owner}/fusion:main-{commit}` - 主分支 + 提交 SHA

## 热修复发布流程

当需要紧急修复生产问题时：

### 1. 创建热修复分支

```bash
# 从需要修复的版本创建分支
git checkout -b hotfix/v1.0.1 v1.0.0
```

### 2. 修复问题

```bash
# 进行修复
git add .
git commit -m "fix: critical bug in authentication"
```

### 3. 合并回 main

```bash
# 合并到 main
git checkout main
git merge hotfix/v1.0.1
git push origin main
```

### 4. 创建修订版本

```bash
git tag -a v1.0.1 -m "Hotfix release v1.0.1

- Fix: Critical authentication bug"
git push origin v1.0.1
```

### 5. 清理热修复分支

```bash
git branch -d hotfix/v1.0.1
```

## 版本回退

如果发布后发现严重问题：

### 1. 删除有问题的标签

```bash
# 删除本地标签
git tag -d v1.0.0

# 删除远程标签
git push origin :refs/tags/v1.0.0
```

### 2. 删除 GitHub Release

在 GitHub Release 页面手动删除该发布。

### 3. 删除 Docker 镜像标签

前往 GitHub Packages 页面删除相应的镜像标签。

### 4. 修复问题后重新发布

修复问题后，可以：
- 使用同样的版本号重新发布（不推荐）
- 增加修订号发布新版本（推荐）

## 常见问题

### Q: 如何修改已发布的 Release 说明？

A: 访问 GitHub Release 页面，点击 "Edit release" 按钮进行编辑。

### Q: 可以删除已发布的版本吗？

A: 可以删除 GitHub Release 和标签，但不建议删除已被用户使用的稳定版本。

### Q: 预发布版本可以升级为稳定版本吗？

A: 不能直接升级。建议删除预发布版本，创建新的稳定版本标签。

### Q: Docker 镜像可以手动推送吗？

A: 可以使用 `make docker-build-multiarch` 和 `make docker-push` 命令手动构建和推送。

### Q: 如何查看所有历史版本？

A: 使用 `git tag -l` 查看所有标签，或访问 GitHub Releases 页面。

## 版本规划建议

### 主版本 (MAJOR)

- 重大架构变更
- 破坏性 API 变更
- 移除已废弃的功能

### 次版本 (MINOR)

- 新功能添加
- 向下兼容的 API 变更
- 标记功能为废弃（但不移除）

### 修订版本 (PATCH)

- Bug 修复
- 安全补丁
- 性能优化
- 文档更新

## 自动化工具

### 查看下一个版本号

```bash
# 查看当前最新标签
git describe --tags --abbrev=0

# 查看所有标签
git tag -l | sort -V
```

### 批量创建标签（不推荐）

```bash
# 创建多个预发布版本
for i in {1..3}; do
  git tag -a "v1.0.0-beta.$i" -m "Beta $i"
done
```

### 验证标签格式

```bash
# 使用正则表达式验证
if [[ $TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+\.[0-9]+)?$ ]]; then
  echo "Valid tag format"
else
  echo "Invalid tag format"
fi
```

## 最佳实践

1. **保持 Changelog 更新**
   - 在每次提交时写清楚的提交信息
   - 定期整理 CHANGELOG.md

2. **使用一致的标签消息**
   - 包含版本号
   - 列出主要变更
   - 提及相关的 Issue 或 PR

3. **充分测试**
   - 在创建标签前运行完整测试套件
   - 使用预发布版本进行集成测试

4. **及时修复问题**
   - 快速响应用户反馈
   - 准备热修复流程

5. **文档先行**
   - 在发布前更新文档
   - 包含迁移指南（如有破坏性变更）

6. **版本兼容性**
   - 维护主版本的稳定性
   - 提前宣布废弃计划

## 相关链接

- [GitHub Releases](https://github.com/{owner}/fusion/releases)
- [GitHub Packages](https://github.com/{owner}/fusion/packages)
- [GitHub Actions](https://github.com/{owner}/fusion/actions)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)