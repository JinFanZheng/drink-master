# 开发工作流程

## 1. 开始新任务

### 前置检查
```bash
git branch --show-current   # 必须在main分支
git status                  # 必须clean
```

### 开始开发
```bash
# 1. 查看任务
gh issue view <issue-id>

# 2. 更新main分支
git checkout main && git pull origin main

# 3. 创建功能分支
git checkout -b feat/<issue-id>-<short-description>

# 4. 标记任务进行中
gh issue edit <issue-id> --add-label "in-progress"
```

## 2. 开发实施

### 使用TodoWrite规划
```
1. 使用TodoWrite工具创建任务计划
2. 分解为具体步骤
3. 实时更新进度 (pending → in_progress → completed)
```

### 质量标准
- ✅ 契约优先：修改`internal/contracts`需记录破坏性变更
- ✅ 类型安全：避免`interface{}`
- ✅ 测试覆盖率≥80%
- ✅ 每次重大修改后运行`make lint && make test`

## 3. 提交代码

### 提交前检查（必须全部通过）
```bash
# 1. 代码格式化
go fmt ./...
find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} \;

# 2. 质量检查
make lint && make test && make build

# 3. 覆盖率验证
go tool cover -func=coverage.out | tail -1  # 必须≥80%
```

### 提交
```bash
git add .
git commit -m "feat: <description>"  # 遵循Conventional Commits
git push -u origin feat/<issue-id>-<short-description>
```

## 4. 创建PR

```bash
gh pr create \
  --title "feat: <description>" \
  --body "## Changes
- Change 1
- Change 2

## Testing
- Test coverage: X%
- All tests pass

Fixes #<issue-id>"
```

## 5. 恢复已有任务

```bash
# 检查当前状态
git branch --show-current   # 应显示feat/<issue-id>-*
git status                  # 查看未提交更改

# 继续开发
gh issue view <issue-id>    # 重新查看需求
# 从上次中断处继续...
```

## 常见问题处理

### 测试覆盖率不足
```bash
# 找出未测试的函数
go tool cover -func=coverage.out | grep "0.0%"
# 为这些函数添加测试
```

### Import格式问题
```bash
# 自动修复import分组
find . -name "*.go" -not -path "./vendor/*" -exec goimports -w {} \;
```

### PR有冲突
```bash
git fetch origin main
git rebase origin/main
# 解决冲突后
git add .
git rebase --continue
git push --force-with-lease
```