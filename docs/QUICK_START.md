# 快速开始指南

## 🚀 最常用命令（一分钟上手）

### 开始任务
```bash
gh issue view <id>                          # 查看任务详情
git checkout main && git pull               # 更新主分支
git checkout -b feat/<id>-description       # 创建功能分支
```

### 提交代码
```bash
make lint && make test && make build        # 质量检查（必须全部通过）
go tool cover -func=coverage.out | tail -1  # 验证覆盖率≥80%
git add . && git commit -m "feat: xxx"      # 提交
git push -u origin feat/<id>-description    # 推送
```

### 创建PR
```bash
gh pr create --title "feat: xxx" --body "Fixes #<id>"
```

## 🔍 遇到问题？

| 问题 | 解决方法 |
|------|----------|
| 测试失败 | 查看 `coverage.out`，补充测试用例 |
| CI失败 | `gh run view` 查看具体原因 |
| 覆盖率不足 | `go tool cover -func=coverage.out \| grep "0.0%"` 找到未测试函数 |
| PR冲突 | `git rebase origin/main` 解决冲突 |
| 需要并行开发 | 参考 `docs/WORKFLOWS/TASK_PARALLEL.md` |
| 紧急修复 | 参考 `docs/EMERGENCY/HOTFIX.md` |

## 📚 更多信息
- 完整开发流程 → `docs/WORKFLOWS/DEVELOPMENT.md`
- PR审查流程 → `docs/WORKFLOWS/PR_MERGE.md`
- 需求分析 → `docs/WORKFLOWS/REQUIREMENT.md`