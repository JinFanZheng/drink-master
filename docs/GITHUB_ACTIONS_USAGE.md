# GitHub Actions 工作流使用指南

本项目包含3个GitHub Actions工作流，各有不同职责，避免重复执行。

## 🚀 工作流概览

### 1. CI (`ci.yml`) - 基础质量检查
**触发条件:** PR创建/更新 + push到main分支  
**职责:** 
- ✅ 代码质量检查 (golangci-lint, gofmt, go vet)
- ✅ 单元测试执行
- ✅ 测试覆盖率验证 (≥78%)  
- ✅ 应用构建验证
- ✅ 启动健康检查

**运行时间:** 每次PR/push都会运行，是其他工作流的基础

### 2. Auto PR Review & Merge (`auto-pr-review.yml`) - 智能合并
**触发条件:** PR状态变更 + CI工作流完成  
**职责:**
- 🔍 等待CI完成并检查结果
- 📊 PR风险等级评估 (低/中/高)
- 🤖 基于风险等级的自动化决策
- ✅ 符合条件的PR自动合并
- 🏷️ 高风险PR添加人工审核标签

**依赖关系:** 依赖CI工作流的成功完成

### 3. GitHub Copilot Integration (`copilot-integration.yml`) - 智能助手
**触发条件:** PR创建/更新 + 特定评论命令  
**职责:**
- 🤖 自动代码审查和建议
- 💬 响应 `@copilot review` 命令
- 💡 响应 `@copilot suggest <需求>` 命令
- 📝 生成智能代码建议

**独立运行:** 不依赖其他工作流

## 📋 使用方式

### 自动触发场景

1. **创建PR时**
   ```
   PR创建 → CI运行 → Auto Review等待CI → (成功后)风险评估 → 自动合并决策
                ↘ Copilot Review生成智能审查报告
   ```

2. **更新PR时**
   ```
   代码推送 → CI重新运行 → Auto Review重新评估 → 更新合并决策
   ```

### 手动触发功能

在PR中评论以下命令：

1. **请求Copilot审查**
   ```
   @copilot review
   ```
   
2. **请求代码建议**
   ```
   @copilot suggest test           # 测试代码建议
   @copilot suggest error handling # 错误处理建议  
   @copilot suggest optimize       # 性能优化建议
   ```

## 🎯 自动合并规则

### 🟢 低风险 - 自动合并
- **条件:** `docs:`、`style:`、`test:`、`chore:` 类型PR
- **要求:** CI通过 + 无冲突 + (可选)Issue链接
- **动作:** 直接squash合并并删除分支

### 🟡 中风险 - 条件合并  
- **条件:** `feat:`、`fix:`、`refactor:` 类型PR
- **要求:** CI通过 + 无冲突 + Issue链接 + 额外验证
- **动作:** 通过后squash合并，否则要求补充信息

### 🔴 高风险 - 人工审核
- **条件:** 契约/数据库/安全相关变更
- **动作:** 添加 `needs-human-review` 标签，等待人工处理

## 📊 状态监控

### 查看工作流状态
```bash
# 列出所有工作流
gh workflow list

# 查看特定工作流运行情况  
gh run list --workflow="CI"
gh run list --workflow="Auto PR Review & Merge"
gh run list --workflow="GitHub Copilot Integration"

# 查看特定运行的详情
gh run view <run-id>
```

### 监控PR处理进度
```bash
# 查看PR检查状态
gh pr checks <pr-number>

# 查看PR详情和标签
gh pr view <pr-number>
```

## 🔧 故障排除

### 常见问题

1. **CI失败导致无法自动合并**
   - 检查CI日志: `gh run view <ci-run-id>`
   - 修复代码质量问题后推送更新

2. **中风险PR要求Issue链接**  
   - 在PR描述中添加 `Fixes #<issue-number>`
   - 或使用 `Closes #<issue-number>`

3. **高风险PR需要人工审核**
   - 查看添加的 `needs-human-review` 标签
   - 联系相关负责人进行审核

4. **Copilot功能无响应**
   - 确认评论格式正确: `@copilot review` 或 `@copilot suggest <需求>`
   - 检查Copilot Integration工作流日志

### 调试命令

```bash
# 手动运行PR验证脚本
./scripts/pr-validation.sh <pr-number>

# 强制重新触发工作流
gh workflow run "Auto PR Review & Merge"

# 查看工作流文件
cat .github/workflows/auto-pr-review.yml
```

## 🎉 最佳实践

1. **PR标题规范**
   - 使用约定式提交格式: `feat:`, `fix:`, `docs:` 等
   - 清晰描述变更内容

2. **PR描述完整**
   - 包含 `Fixes #<issue-id>` 链接
   - 描述变更原因和影响范围

3. **及时响应自动化建议**
   - 关注Copilot审查建议
   - 根据风险等级调整PR内容

4. **监控合并状态**
   - 检查自动合并结果
   - 对失败情况及时处理

## 📞 获取帮助

- 查阅 `docs/AGENT_PR_MERGE_GUIDE.md` 了解详细审查标准
- 查阅 `docs/AUTO_PR_MERGE_SETUP.md` 了解系统配置
- 在团队沟通渠道询问或创建Issue