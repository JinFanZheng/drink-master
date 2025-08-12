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
# 使用worktree创建独立工作目录
git worktree add ../drink-master-<issue-id>-<short-name> -b feat/<issue-id>-<short-name>
cd ../drink-master-<issue-id>-<short-name>          # 切换到worktree目录
make lint && make test && make build  # 基础质量检查
```

**⚠️ 重要约束：**
- **绝对禁止**在非main分支基础上创建新分支
- **必须确保**main分支是最新状态再开始任务
- **必须验证**工作目录干净（无未提交更改）
- **严格遵循**分支命名规范：`feat/<issue-id>-<short-name>`
- **一次只处理一个issue**，禁止并行开发多个任务
- **冲突处理**：如pull时有冲突，必须先完全解决后再继续
- **Worktree管理**：每个任务使用独立的worktree目录，避免分支切换带来的文件变化

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
# 1. 代码格式化和import格式检查（防止CI失败）
go fmt ./...                                                    # 格式化所有Go代码
find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} \;  # 修复import分组
goimports -d $(find . -name "*.go" -not -path "./vendor/*")     # 验证无格式问题（应无输出）

# 2. 质量检查
make lint && make test && make build  # 必须全部通过

# 3. 测试覆盖率检查
go tool cover -func=coverage.out | tail -1  # 必须显示 ≥80.0%

# 4. 提交代码
git add . && git commit -m "feat: ..."  # Conventional Commits 格式
```

**⚠️ 代码质量约束（强制执行）：**
- **测试覆盖率必须达到80%以上**才能提交代码
- 使用 `go tool cover -func=coverage.out` 检查覆盖率
- 如果覆盖率不足80%，**必须**添加更多测试用例
- 重点关注0%覆盖率的函数和方法，优先编写测试
- **代码格式化必须符合Go标准**，使用 `go fmt ./...` 格式化所有代码
- **Import分组格式必须符合项目标准**（三层分组结构）：
  ```go
  import (
      "fmt"           // 第一层：标准库包
      "time"
      
      "github.com/gin-gonic/gin"      // 第二层：外部第三方包  
      "gorm.io/gorm"
      "github.com/shopspring/decimal"
      
      "github.com/ddteam/drink-master/internal/contracts"  // 第三层：项目内部包
      "github.com/ddteam/drink-master/internal/models"
  )
  ```
- **本地goimports处理流程**（防止CI失败）：
  ```bash
  # 1. 安装goimports工具（如未安装）
  go install golang.org/x/tools/cmd/goimports@latest
  
  # 2. 批量修复所有Go文件的import格式
  find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} \;
  
  # 3. 验证格式化是否正确
  goimports -d $(find . -name "*.go" -not -path "./vendor/*")  # 不应有输出
  
  # 4. 提交前最终检查
  make lint && make test && make build  # 必须全部通过
  ```
- **CI goimports失败的常见原因与解决方案**：
  - **原因1**：import分组不正确（缺少空行分隔）
    - **解决**：使用上述goimports命令自动修复
  - **原因2**：存在未使用的import
    - **解决**：goimports会自动移除未使用的import
  - **原因3**：缺少必要的import
    - **解决**：goimports会自动添加缺失的import

### 第5步：PR 创建和清理 (MANDATORY)
```bash
# 在worktree目录中推送分支
git push -u origin feat/<issue-id>-<short-name>
gh pr create --title "..." --body "Fixes #<issue-id> ..."

# PR合并后，清理worktree和本地分支
cd ../drink-master                    # 回到主工作目录
git worktree remove ../drink-master-<issue-id>-<short-name>  # 删除worktree
git branch -d feat/<issue-id>-<short-name>                   # 删除本地分支
git remote prune origin              # 清理远程跟踪分支
```

**违反流程的后果：PR 将被拒绝，需要重新开始。**

### 开发工作流检查清单 ✅
- [ ] **切换主分支并拉取最新代码** (`git checkout main && git pull`)
- [ ] **验证工作目录干净** (`git status`)
- [ ] 查看并理解 Issue 需求
- [ ] 标记 Issue 为 `in-progress` 
- [ ] **创建worktree工作目录**（`git worktree add ../drink-master-<issue-id>-<short-name> -b feat/<issue-id>-<short-name>`）
- [ ] **切换到worktree目录** (`cd ../drink-master-<issue-id>-<short-name>`)
- [ ] 运行基础质量检查
- [ ] **使用 TodoWrite 规划任务**
- [ ] 实施开发并实时更新进度
- [ ] 最终质量检查 (lint + test + build)
- [ ] **验证测试覆盖率 ≥ 80%** (`go tool cover -func=coverage.out | tail -1`)
- [ ] **检查代码格式化** (`go fmt ./...`)
- [ ] **修复import分组格式** (`find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} \;`)
- [ ] **验证import格式正确** (`goimports -d` 命令应无输出)
- [ ] 提交代码并创建 PR
- [ ] 确保 PR 包含 `Fixes #<issue-id>`
- [ ] **PR合并后清理worktree** (`git worktree remove` 和 `git branch -d`)

## 2.1 状态与自动化约定
- 领取：Issue 下评论 `/claim`（自动指派 + 加 `in-progress` 标签 → 看板 In progress）
- 阻塞：加 `blocked` 标签（→ 看板 Blocked），并在评论说明阻塞原因
- 评审：打开 PR 或转为 Ready for review 时，自动为关联 Issue 加 `review` 标签（→ 看板 Review）；若将 PR 转为 Draft 会移除 `review`
- 完成：合并 PR（Fixes #id 自动关闭 Issue）→ 看板 Done

## 2.2 Git Worktree 工作流详解

### Worktree 优势
- **并行开发支持**：每个任务拥有独立的工作目录，无需频繁切换分支
- **文件状态隔离**：避免分支切换时的文件变化影响IDE和构建工具  
- **依赖管理简化**：共享 `.git` 目录，但可独立管理 `node_modules`、`vendor` 等依赖
- **编辑器配置共享**：IDE配置文件共享，提供一致的开发体验

### Worktree 命令参考
```bash
# 创建新的worktree（示例：任务#456）
git worktree add ../drink-master-456-user-auth -b feat/456-user-auth

# 查看所有worktree
git worktree list

# 删除worktree（PR合并后）
git worktree remove ../drink-master-456-user-auth
git branch -d feat/456-user-auth  # 删除本地分支

# 清理无效的worktree引用
git worktree prune

# 修复损坏的worktree（如目录意外删除）
git worktree repair
```

### Worktree 文件结构
```
drink-master/               # 主工作目录（main分支）
├── .git/                  # Git仓库数据（所有worktree共享）
├── docs/
├── internal/
└── ...

drink-master-456-user-auth/ # Task #456的worktree
├── .git -> ../drink-master/.git/worktrees/drink-master-456-user-auth
├── docs/                  # 独立的工作文件
├── internal/
└── ...

drink-master-789-order-api/ # Task #789的worktree  
├── .git -> ../drink-master/.git/worktrees/drink-master-789-order-api
├── docs/
├── internal/
└── ...
```

### 最佳实践
- **命名规范**：严格使用 `../drink-master-<issue-id>-<short-name>` 格式
- **清理时机**：PR合并后立即清理对应的worktree和本地分支
- **共享资源**：IDE配置、Git钩子等自动共享，无需额外配置
- **依赖管理**：每个worktree可以有独立的 `vendor/`、`node_modules/` 等依赖目录

## 2.3 特殊任务类型处理

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
- [ ] **代码格式化检查**：
  - [ ] 运行 `go fmt ./...` 确保无格式化问题
  - [ ] 运行 `find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} \;` 修复import分组
  - [ ] 运行 `goimports -d $(find . -name "*.go" -not -path "./vendor/*")` 验证无格式问题（应无输出）
- [ ] **质量检查**：`make lint && make test && make build` 均通过
- [ ] **测试覆盖率 ≥ 80%**：运行 `go tool cover -func=coverage.out | tail -1` 确认
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

# ✅ 正确：始终基于最新main分支使用worktree
git checkout main && git pull origin main
git worktree add ../drink-master-123-feature-name -b feat/123-feature-name
cd ../drink-master-123-feature-name
```

**场景2：工作目录不干净**
```bash
# ❌ 错误：有未提交更改时切换任务
# 使用worktree时，每个任务有独立目录，但仍需保持干净

# ✅ 正确：在worktree中提交或保存更改
cd ../drink-master-123-current-task
git add . && git commit -m "wip: save current progress"
# worktree的优势：可以直接切换到另一个任务目录而不影响当前工作
cd ../drink-master-456-new-task
```

**场景3：并行开发多个Issue**
```bash
# ❌ 错误：同时在多个分支开发（虽然worktree技术上支持，但不推荐）
git worktree add ../drink-master-111-feature-a -b feat/111-feature-a
git worktree add ../drink-master-222-feature-b -b feat/222-feature-b
# 同时开发两个任务会导致混乱！

# ✅ 正确：完成当前任务后再开始新任务
# 完成feat/111-feature-a，提交PR，合并并清理worktree后再开始222
cd ../drink-master
git worktree remove ../drink-master-111-feature-a
git worktree add ../drink-master-222-feature-b -b feat/222-feature-b
```

### 应急处理
- **Worktree污染**：`cd ../drink-master && git worktree remove ../drink-master-<issue-id>-<name> && git branch -D feat/<issue-id>-<name>` 重新开始
- **提交错误**：在worktree目录中使用 `git reset --soft HEAD~1` 撤销最后一次提交
- **依赖冲突**：先运行 `go mod tidy` 清理依赖后重新构建
- **Worktree目录丢失**：使用 `git worktree prune` 清理无效的worktree引用

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