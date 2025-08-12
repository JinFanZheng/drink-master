# Agent Onboarding Guide

本指南面向新接入的开发代理（agents）。阅读并按本文操作即可在当前进度基础上无缝继续工作。

## 1. 快速路径
- 先读（按序）：
  - `README.md`（工作方式/环境/接口/API文档）
  - `CONTRIBUTING.md`（分支/提交/PR/契约优先）
  - 项目管理：[GitHub Issues](https://github.com/ddteam/drink-master/issues) 和 [项目看板](https://github.com/users/ddteam/projects/1)
- 任务入口（唯一）：
  - Roadmap 看板：`https://github.com/users/ddteam/projects/1`
  - Issues 列表（标签：`backend` / `api` / `docs`，带里程碑）
  - 命令行工具：统一使用 `gh` 操作 GitHub（示例：`gh issue list|view|edit`、`gh pr create`、`gh project item-add`）

## 2. ⚠️ 开发流程（MANDATORY - 强制执行）

**所有开发任务必须严格按照以下流程执行，无例外：**

### 第1步：任务开始前准备 (MANDATORY)
```bash
# 必须按顺序执行：
git checkout main && git pull origin main           # 切换主分支并拉取最新代码
git status                            # 确认工作目录干净
gh issue view <issue-id>              # 查看任务详情
gh issue edit <issue-id> --add-label "in-progress"  # 标记进行中
git checkout -b feat/<issue-id>-<short-name>        # 创建功能分支
make lint && make test && make build  # 基础质量检查
```

**⚠️ 重要约束：**
- **绝对禁止**在非main分支基础上创建新分支
- **必须确保**main分支是最新状态再开始任务
- **必须验证**工作目录干净（无未提交更改）
- **严格遵循**分支命名规范：`feat/<issue-id>-<short-name>`
- **一次只处理一个issue**，禁止并行开发多个任务
- **冲突处理**：如pull时有冲突，必须先完全解决后再继续

### 第2步：任务规划 (MANDATORY)
- **必须使用 TodoWrite 工具**创建详细的任务计划
- **必须**将复杂任务分解为具体步骤
- **必须**实时更新任务进度状态 (pending → in_progress → completed)

### 第3步：开发实施标准
- **契约优先**：如需修改 `internal/contracts/*.go`，必须开 PR 并记录破坏性变更
- **类型安全**：严格使用 Go 类型系统，避免 `interface{}`
- **质量检查**：每次重大修改后运行 `make lint && make test`
- **Mock 开发**：优先使用 `MOCK_MODE=true` 环境变量联调

### 第4步：提交前检查 (MANDATORY)
```bash
make lint && make test && make build  # 必须全部通过
# 检查测试覆盖率是否达到80%+
go tool cover -func=coverage.out | tail -1  # 必须显示 ≥80.0%
# 检查代码格式化
go fmt ./...  # 必须运行以确保代码格式正确
git add . && git commit -m "feat: ..."  # Conventional Commits 格式
```

**⚠️ 代码质量约束（强制执行）：**
- **测试覆盖率必须达到80%以上**才能提交代码
- 使用 `go tool cover -func=coverage.out` 检查覆盖率
- 如果覆盖率不足80%，**必须**添加更多测试用例
- 重点关注0%覆盖率的函数和方法，优先编写测试
- **代码格式化必须符合Go标准**，使用 `go fmt ./...` 格式化所有代码
- CI中的gofmt检查失败会导致构建失败，必须修复后重新提交

### 第5步：PR 创建 (MANDATORY)
```bash
git push -u origin feat/<issue-id>-<short-name>
gh pr create --title "..." --body "Fixes #<issue-id> ..."
```

**违反流程的后果：PR 将被拒绝，需要重新开始。**

### 开发工作流检查清单 ✅
- [ ] **切换主分支并拉取最新代码** (`git checkout main && git pull`)
- [ ] **验证工作目录干净** (`git status`)
- [ ] 查看并理解 Issue 需求
- [ ] 标记 Issue 为 `in-progress` 
- [ ] 创建功能分支（基于最新main分支）
- [ ] 运行基础质量检查
- [ ] **使用 TodoWrite 规划任务**
- [ ] 实施开发并实时更新进度
- [ ] 最终质量检查 (lint + test + build)
- [ ] **验证测试覆盖率 ≥ 80%** (`go tool cover -func=coverage.out | tail -1`)
- [ ] **检查代码格式化** (`go fmt ./...`)
- [ ] 提交代码并创建 PR
- [ ] 确保 PR 包含 `Fixes #<issue-id>`

## 2.1 状态与自动化约定
- 领取：Issue 下评论 `/claim`（自动指派 + 加 `in-progress` 标签 → 看板 In progress）
- 阻塞：加 `blocked` 标签（→ 看板 Blocked），并在评论说明阻塞原因
- 评审：打开 PR 或转为 Ready for review 时，自动为关联 Issue 加 `review` 标签（→ 看板 Review）；若将 PR 转为 Draft 会移除 `review`
- 完成：合并 PR（Fixes #id 自动关闭 Issue）→ 看板 Done

## 2.2 特殊任务类型处理

### Epic 任务处理流程
- **Epic 任务特征**：带 `epic` 标签，包含多个子任务
- **处理策略**：
  1. **不直接开发** Epic 本身，而是协调和跟踪子任务
  2. **检查子任务状态**：使用 `gh issue list` 查看相关子任务进展
  3. **评估整体进度**：确认所有子任务是否完成
  4. **Epic 关闭条件**：所有子任务关闭且验收通过
  5. **文档更新**：确保 Epic 相关文档和指南完整

### 冲突解决标准流程
```bash
# 如果 git pull 时遇到冲突：
git status                              # 查看冲突文件
# 手动解决所有冲突文件
git add .                               # 暂存解决后的文件  
git commit -m "resolve: merge conflicts with main"
# 然后继续正常流程
```

### 并行任务管理约束
- **严格禁止**同时处理多个 issue
- **当前任务未完成前**，不得开始新任务
- **分支切换规则**：只允许在 main 和当前工作分支间切换
- **中断处理**：如需暂停当前任务，必须先提交或 stash 所有更改

## 3. 前后端要点
- 后端：
  - 避免 `interface{}`；遵循 `internal/contracts/*.go` 类型
  - 数据一致性检查（唯一约束、外键约束）
  - 响应时间监控（API响应 < 500ms）
  - 限流（已接入）与必要日志（延迟/告警）
- API设计：
  - RESTful接口设计原则
  - 统一错误处理和响应格式
  - 请求参数验证和边界检查
  - 合理的分页和过滤机制

## 4. 自检清单（提交前）
- [ ] `make lint && make test && make build` 均通过
- [ ] **测试覆盖率 ≥ 80%**：运行 `go tool cover -func=coverage.out | tail -1` 确认
- [ ] **代码格式化正确**：运行 `go fmt ./...` 确保无格式化问题
- [ ] 若改契约：PR 中说明，并更新 README"变更记录"
- [ ] PR 描述包含 `Fixes #<issue-id>`；CI 绿灯
- [ ] 接口返回包含必要信息（如错误码、分页信息），API文档更新
- [ ] **分支基于最新main**：确认分支是从最新main创建
- [ ] **无未追踪文件**：`git status` 显示工作目录干净
- [ ] **提交信息规范**：遵循 Conventional Commits 格式

## 5. 验收标准（MVP DoD 摘要）
- API响应时间 < 500ms，错误处理完整
- 数据库事务一致性，支持并发访问
- **单元测试覆盖率 ≥ 80%**，集成测试通过
- `/api/health`、`/api/drinks/*` 均可用且有完整错误码
- 文档到位；Mock 支持可独立联调；CI 绿灯

**测试覆盖率验收细则：**
- 使用 `make test` 生成覆盖率报告
- 总体覆盖率必须达到80.0%或以上
- 新增代码的覆盖率应达到90%以上
- 所有核心业务逻辑函数必须有相应测试
- 错误处理分支也需要测试覆盖

## 6. 常见问题与故障排除

### 常见错误场景
**场景1：基于错误分支创建新分支**
```bash
# ❌ 错误：在feature分支上创建新分支
git checkout feat/old-branch
git checkout -b feat/new-branch  # 这是错误的！

# ✅ 正确：始终基于最新main分支
git checkout main && git pull origin main
git checkout -b feat/new-branch
```

**场景2：工作目录不干净**
```bash
# ❌ 错误：有未提交更改时切换任务
# 未提交的文件会污染新分支

# ✅ 正确：先清理工作目录
git add . && git commit -m "wip: save current progress"
# 或者
git stash push -m "临时保存更改"
```

**场景3：并行开发多个Issue**
```bash
# ❌ 错误：同时在多个分支开发
git checkout feat/issue-1
# 开发一半，切换到另一个issue
git checkout -b feat/issue-2  # 这会导致混乱！

# ✅ 正确：完成当前任务后再开始新任务
# 完成feat/issue-1，提交PR，合并后再开始issue-2
```

### 应急处理
- **分支污染**：`git checkout main && git branch -D <polluted-branch>` 重新开始
- **提交错误**：使用 `git reset --soft HEAD~1` 撤销最后一次提交
- **依赖冲突**：先运行 `go mod tidy` 清理依赖后重新构建

## 7. 关键链接
- 仓库：`https://github.com/ddteam/drink-master`
- 看板：`https://github.com/users/ddteam/projects/1`
- 文档入口：`README.md`（包含完整开发指南）、`CONTRIBUTING.md`、`docs/AGENT_PR_MERGE_GUIDE.md`

## 8. 流程合规性检查

在开始任何新任务前，**必须**完成以下合规性检查：

```bash
# 完整的任务启动检查脚本
echo "=== 任务启动合规性检查 ==="
echo "1. 检查当前分支..."
git branch --show-current  # 应该显示 main

echo "2. 检查工作目录状态..."
git status  # 应该显示 "working tree clean"

echo "3. 更新本地主分支..."
git pull origin main  # 应该显示 "Already up to date" 或成功更新

echo "4. 确认无本地未推送分支..."
git branch --no-merged main  # 不应该有输出

echo "✅ 如果以上检查全部通过，可以开始新任务"
```

**只有通过所有检查，才能继续执行标准开发流程！**