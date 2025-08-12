# Agent PR Merge Guide

本指南面向AI开发代理（agents），提供安全、标准化的Pull Request合并操作流程。

## 核心原则

- **安全第一**: 充分验证后再合并，避免破坏主分支
- **使用 gh 命令**: 统一使用 GitHub CLI 进行所有 GitHub 相关操作
- **自动化优先**: 能自动检查的尽量自动化，减少人工介入
- **清晰沟通**: 合并前后及时通知，问题及时上报

## 1. 准备工作

### 查看待合并的PR
```bash
# 列出所有开放的PR
gh pr list --state open

# 查看特定PR详情
gh pr view <pr-number>

# 查看PR状态检查
gh pr checks <pr-number>
```

### 了解变更内容
```bash
# 查看PR的文件变更
gh pr diff <pr-number>

# 查看PR关联的Issue
gh pr view <pr-number> --json body,title | jq -r '.body' | grep -E '(Fixes|Closes) #[0-9]+'
```

## 2. 合并前检查清单

### 必需检查项 (❌ 任一项不通过则不可合并)

```bash
# 1. CI/CD状态检查
gh pr checks <pr-number>
# 确保所有检查都是 ✓ PASS 状态

# 2. 冲突检查
gh pr view <pr-number> --json mergeable
# 确保 mergeable: true

# 3. Issue关联检查
gh pr view <pr-number> --json body | jq -r '.body' | grep -E '(Fixes|Closes) #[0-9]+'
# 确保有正确的Issue链接

# 4. 代码质量门槛检查
gh pr checks <pr-number> | grep -E '(lint|test|coverage)'
# 确保以下检查通过：
# - golangci-lint 代码质量检查
# - 所有单元测试通过
# - 测试覆盖率 ≥78%
# - 循环复杂度 ≤15 (gocyclo)
# - 代码格式检查 (gofmt)
```

### 功能完整性验证

基于项目的测试计划和验收标准，agents需要验证系统功能的完整性：

```bash
# 4. 健康检查 - 基础系统状态
curl -s http://localhost:8080/api/health | jq '.'
# 确保返回 {"status": "ok", ...}

# 5. 数据库连接测试 - 核心数据层验证
curl -s http://localhost:8080/api/health/db | jq '.'
# 确保数据库连接正常

# 6. 基础API功能测试 - 核心业务逻辑验证
curl -s -X GET http://localhost:8080/api/drinks | jq '.data | length'
# 确保API响应正常格式

# 7. CRUD操作测试（如果涉及数据操作变更）
curl -s -X POST http://localhost:8080/api/drinks \
  -H "Content-Type: application/json" \
  -d '{"name":"测试饮品","category":"coffee","price":25.5}' \
  | jq '.data.id'
# 确保创建操作正常
```

### 业务逻辑验证

针对关键业务功能的验证：

```bash
# 8. 验证契约一致性（如果涉及 internal/contracts 变更）
# 检查API契约是否同步
gh pr diff <pr-number> --name-only | grep -E 'internal/contracts' && echo "⚠️ 契约变更需要验证API一致性"

# 9. 验证数据模型完整性（如果涉及数据库变更）
# 检查模型定义和数据库迁移是否一致
gh pr diff <pr-number> --name-only | grep -E '(models|migrations)' && echo "⚠️ 数据模型变更需要验证"

# 10. 验证认证授权功能（如果涉及安全相关变更）
# 检查JWT和权限验证是否正常
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}' \
  | jq '.token'
```

### 变更影响评估
```bash
# 查看变更的文件类型和范围
gh pr diff <pr-number> --name-only

# 检查是否有契约变更（需特别注意）
gh pr diff <pr-number> --name-only | grep -E '(internal/contracts|models)'

# 检查是否有breaking changes
gh pr view <pr-number> --json body,title | grep -i -E '(breaking|BREAKING)'
```

## 3. 风险分类与合并策略

### 🟢 低风险 - 可自动合并
满足以下条件的PR可以直接合并：
- CI/CD全部通过 ✓
- 无merge conflicts ✓  
- 有正确的Issue链接 ✓
- 属于以下类型之一：
  - `docs:` 文档更新
  - `style:` 样式调整，不影响逻辑
  - `test:` 测试用例添加/修复
  - `chore:` 工具配置、依赖更新（非breaking）

```bash
# 自动合并命令
gh pr merge <pr-number> --merge --delete-branch
```

### 🟡 中风险 - 需额外验证
以下类型需要更仔细的检查：
- `feat:` 新功能实现
- `fix:` Bug修复
- `refactor:` 代码重构

**额外检查步骤：**
```bash
# 1. 运行完整的质量检查
make lint && make test && make build
# 注意：确保测试覆盖率达到78%以上

# 2. 代码复杂度检查
golangci-lint run --timeout=5m
# 重点检查循环复杂度是否超过15
# 如发现复杂度过高，建议重构为多个小函数

# 3. 功能回归测试
make health-check  # 健康状态检查
make test-api      # API功能测试

# 3. 查看PR大小（行数变更）
gh pr diff <pr-number> --stat

# 4. 检查是否影响关键API接口
gh pr diff <pr-number> --name-only | grep -E '(internal/handlers|internal/contracts)'

# 5. 验证当前里程碑的验收标准（如果是功能PR）
# 检查PR关联的Issue所属里程碑，验证相应功能要求
MILESTONE=$(gh issue view $(gh pr view <pr-number> --json body | grep -oE '#[0-9]+' | head -1 | cut -c2-) --json milestone | jq -r '.milestone.title // "无里程碑"')
echo "验证里程碑: $MILESTONE 的功能要求"

# 核心业务逻辑验证（适用所有功能PR）：
# - API健康状态正常
# - 核心功能响应正确格式
# - 数据库操作正常
# - 无阻塞性错误或异常

# 6. 数据库迁移验证（如果涉及数据库变更）
MIGRATION_CHANGES=$(gh pr diff <pr-number> --name-only | grep -E '(migrations/|models/)')
if [ ! -z "$MIGRATION_CHANGES" ]; then
  echo "检测到数据库变更，验证迁移脚本..."
  # 确保迁移可以正常执行且可回滚
fi
```

**合并条件：**
- `make lint && make test && make build` 全部通过
- **测试覆盖率 ≥78%** (CI自动检查)
- **代码复杂度合规** (单个函数循环复杂度 ≤15)
- **代码格式规范** (通过gofmt和golangci-lint检查)
- 核心API功能测试通过（/api/health, /api/drinks/*）
- 变更行数 < 500行 或 新增功能完整且独立
- 无API breaking changes
- 符合相应的MVP验收标准

```bash
# 合并命令（优先使用squash）
gh pr merge <pr-number> --squash --delete-branch
```

### 🔴 高风险 - 必须人工审核
以下情况**不可自动合并**，需要人工介入：
- 契约变更 (`internal/contracts/*`)
- 数据库schema变更 (`migrations/*`, `internal/models/*`)
- 安全相关修改（认证、权限、加密）
- 配置文件大幅变更
- 跨多个模块的大规模重构
- 有 `BREAKING CHANGE` 标记

**处理方式：**
```bash
# 添加需要人工审核的标签
gh pr edit <pr-number> --add-label "needs-human-review"

# 请求特定人员审核
gh pr edit <pr-number> --add-reviewer <reviewer-username>

# 添加评论说明风险点
gh pr comment <pr-number> --body "⚠️ 此PR包含高风险变更，已标记需要人工审核：
- [具体风险描述]
- [影响范围说明]
请相关负责人审核后手动合并。"
```

## 4. 特殊情况处理

### 合并冲突解决
```bash
# 1. 检查冲突详情
gh pr view <pr-number> --json mergeable,mergeStateStatus

# 2. 如果是简单的自动可解决冲突
gh pr comment <pr-number> --body "发现合并冲突，请作者更新分支：\`git merge main\` 或 \`git rebase main\`"

# 3. 建议作者更新分支
gh pr edit <pr-number> --add-label "needs-rebase"
```

### 紧急修复流程
对于标记为 `urgent` 或 `hotfix` 的PR：
```bash
# 1. 快速验证基本检查
gh pr checks <pr-number> | head -5

# 2. 直接合并（跳过某些非关键检查）
gh pr merge <pr-number> --merge --delete-branch

# 3. 立即通知
gh pr comment <pr-number> --body "🚨 紧急修复已合并并部署。请相关人员关注生产环境状态。"
```

## 5. 合并后操作

### 自动化任务
```bash
# 1. 更新Issue状态（通过Fixes #xx自动关闭）

# 2. 通知相关人员
gh pr view <pr-number> --json author,assignees

# 3. 检查部署状态（如果有自动部署）
gh run list --limit 1 --workflow=deploy
```

### 问题上报
如果合并后发现问题：
```bash
# 1. 创建回滚Issue
gh issue create --title "回滚PR #<pr-number>: [问题描述]" \
  --body "PR #<pr-number> 合并后发现问题，需要紧急回滚。

**问题描述**: [具体问题]
**影响范围**: [影响范围]
**回滚方案**: [回滚步骤]

原PR: #<pr-number>" \
  --label "urgent,rollback"

# 2. 如果需要立即回滚
git revert <commit-hash> --no-edit
git push origin main
```

## 6. 常用命令参考

### PR查看和管理
```bash
# 查看PR列表
gh pr list --limit 10 --state open

# 查看PR详情
gh pr view <pr-number>

# 查看PR检查状态
gh pr checks <pr-number>

# 查看PR变更
gh pr diff <pr-number>

# 添加标签
gh pr edit <pr-number> --add-label "label-name"

# 添加评论
gh pr comment <pr-number> --body "评论内容"

# 请求审核
gh pr edit <pr-number> --add-reviewer "username"
```

### PR合并
```bash
# Merge commit (保留提交历史)
gh pr merge <pr-number> --merge --delete-branch

# Squash merge (压缩为单个提交，推荐)
gh pr merge <pr-number> --squash --delete-branch

# Rebase merge (变基合并)
gh pr merge <pr-number> --rebase --delete-branch
```

### Issue管理
```bash
# 查看Issue详情
gh issue view <issue-number>

# 更新Issue状态
gh issue edit <issue-number> --add-label "label-name"

# 创建新Issue
gh issue create --title "标题" --body "内容" --label "标签"
```

## 7. 决策流程图

```
PR待合并
    ↓
CI/CD是否全通过？
    ↓ NO → 等待修复 → 通知作者
    ↓ YES
是否有冲突？
    ↓ YES → 通知作者解决冲突
    ↓ NO
是否有Issue链接？
    ↓ NO → 添加评论要求补充
    ↓ YES
变更类型？
    ↓
docs/style/test → 自动合并
feat/fix/refactor → 额外验证 → 满足条件？→ 合并
contracts/db/security → 人工审核标记
```

## 8. 自动化验证脚本

基于Go项目的测试指南，以下是一个完整的PR验证脚本示例：

```bash
#!/bin/bash
# pr-validation.sh - PR合并前自动验证脚本

PR_NUMBER=$1
if [ -z "$PR_NUMBER" ]; then
  echo "Usage: $0 <pr-number>"
  exit 1
fi

echo "🔍 开始验证PR #$PR_NUMBER..."

# 1. 基础检查
echo "1️⃣ 检查CI状态..."
if ! gh pr checks $PR_NUMBER | grep -q "✓"; then
  echo "❌ CI检查未通过"
  exit 1
fi

# 2. 检查可合并状态
echo "2️⃣ 检查合并状态..."
MERGEABLE=$(gh pr view $PR_NUMBER --json mergeable | jq -r '.mergeable')
if [ "$MERGEABLE" != "true" ]; then
  echo "❌ PR有冲突，无法合并"
  exit 1
fi

# 3. 检查Issue链接
echo "3️⃣ 检查Issue关联..."
if ! gh pr view $PR_NUMBER --json body | jq -r '.body' | grep -E '(Fixes|Closes) #[0-9]+'; then
  echo "⚠️ 未发现Issue链接，请确认"
fi

# 4. 代码质量检查
echo "4️⃣ 运行代码质量检查..."
if ! make lint > /dev/null 2>&1; then
  echo "❌ 代码质量检查失败"
  exit 1
fi

# 4.1 检查测试覆盖率
echo "4.1️⃣ 检查测试覆盖率..."
COVERAGE=$(go test -coverprofile=coverage.out ./... > /dev/null 2>&1 && go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
if (( $(echo "$COVERAGE < 78" | bc -l) )); then
  echo "❌ 测试覆盖率不足: ${COVERAGE}% (要求≥78%)"
  exit 1
fi
echo "✅ 测试覆盖率: ${COVERAGE}%"

# 4.2 检查循环复杂度
echo "4.2️⃣ 检查循环复杂度..."
if golangci-lint run --disable-all --enable=gocyclo > /dev/null 2>&1; then
  echo "✅ 循环复杂度检查通过"
else
  echo "❌ 发现循环复杂度过高的函数，请重构"
  exit 1
fi

# 5. 功能验证
echo "5️⃣ 运行功能测试..."
if ! make health-check > /dev/null 2>&1; then
  echo "❌ 健康检查失败"
  exit 1
fi

# 6. API测试
echo "6️⃣ 验证API功能..."
if ! make test-api > /dev/null 2>&1; then
  echo "❌ API测试失败"
  exit 1
fi

# 7. 检查变更类型、里程碑和风险
echo "7️⃣ 分析变更风险和里程碑验证..."
CHANGED_FILES=$(gh pr diff $PR_NUMBER --name-only)
RISK_LEVEL="low"

# 检查关联Issue的里程碑
ISSUE_NUM=$(gh pr view $PR_NUMBER --json body | jq -r '.body' | grep -oE 'Fixes #[0-9]+' | grep -oE '[0-9]+' | head -1)
if [ ! -z "$ISSUE_NUM" ]; then
  MILESTONE=$(gh issue view $ISSUE_NUM --json milestone | jq -r '.milestone.title // "无里程碑"')
  echo "📋 关联Issue #$ISSUE_NUM，里程碑: $MILESTONE"
fi

# 风险评估
if echo "$CHANGED_FILES" | grep -E '(internal/contracts|internal/handlers)'; then
  RISK_LEVEL="medium"
  echo "⚠️ 发现API或契约变更，风险等级: 中"
fi

if echo "$CHANGED_FILES" | grep -E '(migrations/|internal/models|security|auth)'; then
  RISK_LEVEL="high"
  echo "🚨 发现高风险变更（数据库/安全），需要人工审核"
  gh pr edit $PR_NUMBER --add-label "needs-human-review"
  exit 2
fi

echo "✅ PR #$PR_NUMBER 验证通过，风险等级: $RISK_LEVEL"
echo "可以安全合并"
```

## 9. 测试质量专项检查

基于项目的测试覆盖率和代码质量要求，以下检查是必须的：

### 测试覆盖率检查
```bash
# 检查当前覆盖率
COVERAGE=$(go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "当前测试覆盖率: ${COVERAGE}%"

# 验证是否达到最低要求
if (( $(echo "$COVERAGE < 78" | bc -l) )); then
  echo "❌ 测试覆盖率不足: ${COVERAGE}% (要求≥78%)"
  echo "请添加更多测试用例或联系团队调整阈值"
  exit 1
else
  echo "✅ 测试覆盖率符合要求: ${COVERAGE}%"
fi
```

### 代码复杂度检查
```bash
# 检查循环复杂度 (gocyclo)
# 限制: 单个函数复杂度不能超过15
echo "检查循环复杂度..."
if golangci-lint run --disable-all --enable=gocyclo; then
  echo "✅ 循环复杂度检查通过"
else
  echo "❌ 发现复杂度过高的函数"
  echo "解决方案：将大函数拆分为多个小函数，或使用table-driven tests"
  exit 1
fi
```

### Pre-commit Hook检查
```bash
# 验证是否设置pre-commit hooks
if [ -f ".githooks/pre-commit" ] && [ -x ".githooks/pre-commit" ]; then
  echo "✅ Pre-commit hook已设置"
  # 检查Git hooks路径配置
  HOOKS_PATH=$(git config core.hooksPath)
  if [ "$HOOKS_PATH" = ".githooks" ]; then
    echo "✅ Git hooks路径配置正确"
  else
    echo "⚠️ 建议设置: git config core.hooksPath .githooks"
  fi
else
  echo "⚠️ 未检测到pre-commit hook，建议使用项目提供的.githooks/pre-commit"
fi
```

### 测试类型覆盖检查
```bash
# 检查各模块的测试覆盖情况
go test -coverprofile=coverage.out ./...
echo "各模块覆盖率详情："
go tool cover -func=coverage.out | grep -E "(internal/handlers|internal/services|internal/repositories|pkg/)"

# 重点关注低覆盖率模块
LOW_COVERAGE=$(go tool cover -func=coverage.out | awk '$3 ~ /%/ && $3 < "70.0%" {print $1 ": " $3}' | head -5)
if [ ! -z "$LOW_COVERAGE" ]; then
  echo "⚠️ 以下模块覆盖率较低（<70%）："
  echo "$LOW_COVERAGE"
  echo "请在后续开发中优先补充这些模块的测试"
fi
```

## 10. 最佳实践

1. **使用验证脚本**: 运行上述脚本进行全面检查
2. **遵循现有工具链**: 优先使用项目的 `make` 命令而不是直接的go命令  
3. **Pre-commit检查**: 确保本地开发环境使用pre-commit hooks
4. **测试覆盖率监控**: 定期检查各模块覆盖率，优先补充低覆盖模块
5. **代码复杂度控制**: 及时重构过于复杂的函数，保持代码可读性
6. **批量处理**: 使用脚本批量检查多个PR状态
7. **定期清理**: 定期检查长时间未更新的PR
8. **监控部署**: 合并后关注部署状态和错误日志
9. **文档更新**: 重要变更及时更新相关文档
10. **团队协作**: 重要决策及时在团队群组通知
11. **Mock优先**: 利用项目的Mock模式进行快速验证

## 11. 应急联系

如果遇到无法自动处理的复杂情况：
1. 在相关Issue或PR中 @mention 项目负责人
2. 添加 `needs-human-review` 标签
3. 在团队沟通渠道中报告情况

---

**注意**: 此指南会根据项目发展持续更新，请agents定期检查最新版本。