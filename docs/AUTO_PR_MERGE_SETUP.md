# 自动PR审查和合并系统设置指南

本文档描述了基于 `docs/AGENT_PR_MERGE_GUIDE.md` 实现的自动化PR审查和合并系统。

## 🏗️ 系统架构

### 核心组件

1. **GitHub Actions 工作流**
   - `auto-pr-review.yml` - 主要的自动审查和合并流程
   - `copilot-integration.yml` - GitHub Copilot 集成功能

2. **验证脚本**
   - `scripts/pr-validation.sh` - 独立的PR验证脚本

3. **风险分级系统**
   - 🟢 **低风险**: 自动合并
   - 🟡 **中风险**: 条件合并
   - 🔴 **高风险**: 人工审核

## 🔧 设置步骤

### 1. 权限配置

确保GitHub Actions具有必要权限：

```yaml
permissions:
  contents: write
  pull-requests: write
  issues: write
```

### 2. 环境变量配置

在GitHub仓库设置中添加以下secrets（如需要）：

- `GITHUB_TOKEN` (自动提供)
- `COPILOT_API_KEY` (如果使用真实的Copilot API)

### 3. 分支保护规则

建议配置以下分支保护规则：

```json
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "CI",
      "Auto PR Review & Merge"
    ]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": null,
  "restrictions": null,
  "allow_auto_merge": true
}
```

## 🚀 工作流程说明

### 自动PR审查流程

1. **触发条件**
   - PR创建、更新或标记为ready for review
   - 目标分支为main

2. **执行步骤**
   ```mermaid
   graph TD
       A[PR事件触发] --> B[风险评估]
       B --> C{风险等级}
       C -->|低| D[基础检查]
       C -->|中| E[扩展验证]
       C -->|高| F[人工审核标记]
       D --> G[自动合并]
       E --> H[条件合并]
       F --> I[等待人工处理]
   ```

3. **质量门槛**
   - CI/CD全部通过
   - 无合并冲突
   - 测试覆盖率 ≥78%
   - 代码复杂度符合要求

### GitHub Copilot集成

1. **代码审查**
   - 自动触发：PR创建/更新时
   - 手动触发：在PR中评论 `@copilot review`

2. **代码建议**
   - 评论触发：`@copilot suggest <请求内容>`
   - 支持测试、错误处理、优化等建议

## 📊 风险分级详解

### 🟢 低风险 - 自动合并
**条件:**
- PR类型：`docs:`、`style:`、`test:`、`chore:`
- CI全部通过
- 无冲突
- 有Issue链接（可选）

**处理:**
```bash
gh pr merge $PR_NUMBER --squash --delete-branch
```

### 🟡 中风险 - 条件合并
**条件:**
- PR类型：`feat:`、`fix:`、`refactor:`
- 变更文件涉及业务逻辑
- 所有质量检查通过

**额外验证:**
- 功能测试
- API健康检查
- Copilot代码审查

### 🔴 高风险 - 人工审核
**触发条件:**
- 契约文件变更 (`internal/contracts/*`)
- 数据库变更 (`migrations/*`, `models/*`)
- 安全相关 (`auth`, `security`)
- 有 `BREAKING CHANGE` 标记

**处理:**
- 添加 `needs-human-review` 标签
- 请求相关人员审核
- 不执行自动合并

## 🧪 测试和验证

### 本地测试

1. **验证脚本测试**
   ```bash
   # 测试验证脚本（需要有效的PR编号）
   ./scripts/pr-validation.sh 123
   ```

2. **模拟GitHub Actions**
   ```bash
   # 使用act工具本地运行
   act pull_request -j auto-review
   ```

### 功能测试清单

- [ ] 低风险PR自动合并
- [ ] 中风险PR条件合并
- [ ] 高风险PR人工标记
- [ ] Copilot审查生成
- [ ] 错误情况处理

## 📈 监控和维护

### 监控指标

1. **合并成功率**
   - 自动合并成功的PR数量
   - 失败需要人工干预的情况

2. **质量指标**
   - 测试覆盖率趋势
   - 代码复杂度分布

3. **效率指标**
   - PR处理时间
   - 人工干预频率

### 日志检查

查看GitHub Actions运行日志：
```bash
gh run list --workflow="Auto PR Review & Merge"
gh run view <run-id>
```

## 🔧 故障排除

### 常见问题

1. **权限不足**
   - 确认GitHub Actions权限配置
   - 检查GITHUB_TOKEN是否有效

2. **合并冲突**
   - 自动添加 `needs-rebase` 标签
   - 通知作者解决冲突

3. **测试失败**
   - 检查CI配置
   - 确认测试环境设置

### 调试步骤

1. **检查工作流状态**
   ```bash
   gh workflow list
   gh workflow view "Auto PR Review & Merge"
   ```

2. **查看具体运行**
   ```bash
   gh run list --workflow="auto-pr-review.yml"
   ```

3. **手动验证**
   ```bash
   ./scripts/pr-validation.sh <pr-number>
   ```

## 🚀 升级和扩展

### 计划增强功能

1. **真实Copilot API集成**
2. **更智能的风险评估**
3. **自定义审查规则**
4. **团队通知集成**

### 配置自定义

可以通过修改以下文件自定义行为：
- `.github/workflows/auto-pr-review.yml` - 主流程
- `scripts/pr-validation.sh` - 验证逻辑
- `docs/AGENT_PR_MERGE_GUIDE.md` - 审查标准

## 📞 支持和反馈

如遇到问题或需要功能改进：

1. 创建Issue并添加 `automation` 标签
2. 在团队沟通渠道报告
3. 查阅 `docs/AGENT_PR_MERGE_GUIDE.md` 获取详细标准

---

**注意**: 此系统持续演进中，请定期检查更新和最佳实践。