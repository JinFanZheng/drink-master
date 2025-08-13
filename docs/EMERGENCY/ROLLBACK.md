# 回滚操作指南

## 🔄 代码回滚

### 1. 回滚最近的提交
```bash
# 查看最近的提交
git log --oneline -5

# 回滚最后一次提交（保留更改）
git reset --soft HEAD~1

# 回滚最后一次提交（丢弃更改）
git reset --hard HEAD~1
```

### 2. 回滚已合并的PR
```bash
# 方法1：使用GitHub CLI
gh pr view <pr-number>  # 查看PR信息
git revert -m 1 <merge-commit-hash>
git push origin main

# 方法2：创建revert PR
git checkout main && git pull
git checkout -b revert-<pr-number>
git revert <commit-hash>
git push -u origin revert-<pr-number>
gh pr create --title "Revert PR #<pr-number>" --label "revert"
```

### 3. 回滚到指定版本
```bash
# 查找目标版本
git log --oneline --graph

# 回滚到指定commit
git checkout main
git reset --hard <commit-hash>
git push --force-with-lease origin main  # 谨慎使用！
```

## 🗄️ 数据库回滚

### 1. 回滚最近的迁移
```bash
# 如果使用migrate工具
make db-rollback

# 手动回滚
mysql -u root -p vending_machine_dev < migrations/rollback/<version>.sql
```

### 2. 数据恢复
```sql
-- 从备份恢复
SOURCE /backup/backup_20240101.sql;

-- 或使用事务日志恢复
-- 根据具体数据库配置
```

## 📦 依赖回滚

### Go依赖回滚
```bash
# 查看依赖历史
git log go.mod

# 恢复到之前的版本
git checkout <commit-hash> -- go.mod go.sum
go mod download
go mod tidy
```

## 🚨 紧急回滚检查清单

### 回滚前确认
- [ ] 确定回滚范围（代码/数据库/配置）
- [ ] 评估回滚影响
- [ ] 通知相关人员
- [ ] 准备回滚方案

### 回滚执行
- [ ] 执行回滚操作
- [ ] 验证回滚结果
- [ ] 监控系统状态
- [ ] 更新Issue状态

### 回滚后跟进
- [ ] 创建问题分析Issue
- [ ] 记录回滚原因
- [ ] 制定修复计划
- [ ] 更新文档

## 常见回滚场景

### 场景1：API破坏性变更
```bash
# 快速回滚API版本
git revert <api-change-commit>
# 保持向后兼容
```

### 场景2：性能严重下降
```bash
# 回滚到性能正常的版本
git bisect start
git bisect bad  # 当前版本性能差
git bisect good <known-good-commit>
# 找到问题commit后回滚
```

### 场景3：数据损坏
```sql
-- 使用备份恢复特定表
DROP TABLE IF EXISTS corrupted_table;
CREATE TABLE corrupted_table LIKE backup.corrupted_table;
INSERT INTO corrupted_table SELECT * FROM backup.corrupted_table;
```

## 预防措施

### 1. 特性开关
```go
if feature.IsEnabled("new-feature") {
    // 新功能代码
} else {
    // 旧功能代码
}
```

### 2. 灰度发布
- 先发布到测试环境
- 小流量验证
- 逐步放量

### 3. 自动回滚
```yaml
# CI/CD配置
rollback:
  triggers:
    - error_rate > 5%
    - response_time > 2s
    - health_check_fail
```