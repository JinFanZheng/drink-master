# 紧急修复流程

## 🚨 生产Bug紧急修复

### 1. 创建hotfix分支
```bash
# 基于main分支（不是feature分支！）
git checkout main && git pull origin main
git checkout -b hotfix/<issue-description>
```

### 2. 快速修复
```bash
# 修改代码...
# 最小化改动，只修复紧急问题

# 快速测试（至少保证不破坏现有功能）
make test  # 可以暂时忽略覆盖率要求
```

### 3. 提交和推送
```bash
git add .
git commit -m "hotfix: <紧急问题描述>"
git push -u origin hotfix/<issue-description>
```

### 4. 创建紧急PR
```bash
gh pr create \
  --title "🚨 HOTFIX: <问题描述>" \
  --label "urgent,hotfix" \
  --body "## 紧急问题
- 问题描述：XXX
- 影响范围：XXX
- 修复方案：XXX

## 测试
- [ ] 本地测试通过
- [ ] 不破坏现有功能

**需要立即合并！**"
```

### 5. 快速合并
```bash
# 通知相关人员review
# 如果CI基础检查通过，可以直接合并
gh pr merge <pr-number> --merge  # 使用merge保留历史
```

## 🔄 快速回滚

### 如果hotfix引入新问题
```bash
# 方法1：revert commit
git checkout main && git pull
git revert <commit-hash> --no-edit
git push origin main

# 方法2：创建回滚PR
gh pr create \
  --title "🔄 Revert: <原PR标题>" \
  --label "urgent,revert" \
  --body "Reverting #<original-pr-number> due to <reason>"
```

## ⚡ 跳过检查的情况

以下情况可以跳过某些检查：
- ✅ 可跳过80%覆盖率要求
- ✅ 可跳过完整的lint检查
- ❌ 不能跳过编译检查
- ❌ 不能跳过基础测试

## 📋 Hotfix后的后续工作

1. **创建正式修复Issue**
```bash
gh issue create \
  --title "[Tech Debt] 完善<hotfix内容>的正式修复" \
  --label "tech-debt,priority-high" \
  --body "Hotfix PR: #<hotfix-pr-number>
需要：
- [ ] 补充测试用例
- [ ] 完善错误处理
- [ ] 代码优化"
```

2. **更新文档**（如果需要）

3. **根因分析**（事后进行）

## 常见紧急场景

### 场景1：API返回错误
```go
// 快速修复：添加错误处理
if err != nil {
    // hotfix: 临时返回友好错误
    c.JSON(500, gin.H{"error": "服务暂时不可用"})
    return
}
```

### 场景2：数据库连接失败
```go
// 快速修复：添加重试逻辑
for i := 0; i < 3; i++ {
    if err := db.Ping(); err == nil {
        break
    }
    time.Sleep(time.Second)
}
```

### 场景3：第三方服务异常
```go
// 快速修复：降级处理
if wechatErr != nil {
    // 暂时跳过微信登录，使用备用方案
    return backupLogin()
}
```