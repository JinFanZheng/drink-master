# Claude Code Hooks 配置指南

本项目使用 Claude Code hooks 来强制执行开发流程和验收标准，确保代码质量和流程合规性。

## 🎯 解决的问题

**主要问题**: Claude Code 在长时间任务后容易"忘记"执行验收标准和后续步骤。

**解决方案**: 通过 hooks 系统自动检查和强制执行开发流程的每个环节。

## 📋 Hook 配置概述

### Hook 事件类型

1. **PostToolUse Hook** - 工具执行后检查
   - 在每次 Edit、Write、Bash 等工具使用后自动触发
   - 强制执行开发验收标准
   - 阻断不合规的操作

2. **PreToolUse Hook** - 工具执行前验证  
   - PR创建前验证所有前置条件
   - 提交前验证代码质量和测试覆盖率
   - 防止跳过必要步骤

3. **UserPromptSubmit Hook** - 用户输入验证
   - 检测"完成"、"done"等关键词
   - 强制验证任务是否真正完成所有验收标准

## 🚀 使用方法

### 1. 启用 Hooks

将 hooks 配置添加到 Claude Code 设置中：

```bash
# 方法1: 复制配置到用户设置
cp .claude/settings.json ~/.claude/settings.json

# 方法2: 在项目根目录使用（推荐）
# Claude Code 会自动读取项目根目录的 .claude/settings.json
```

### 2. Hook 工作流程

#### 开发过程中的自动检查
```bash
# 每次代码编辑后，会自动检查：
# ✅ 代码格式化 (go fmt)
# ✅ Import 格式化 (goimports) 
# ✅ 构建状态 (make build)
# ✅ 测试状态 (make test)
# ✅ Lint 检查 (make lint)
# ✅ 测试覆盖率 ≥ 80%

# 如果检查失败，会阻断后续操作并提供修复建议
```

#### PR 创建前的强制验证
```bash
# 当 Claude 尝试执行 "gh pr create" 时：
# ✅ 验证在正确的 feature 分支
# ✅ 验证工作目录干净
# ✅ 验证所有验收检查通过
# ✅ 最终质量检查

# 只有全部通过才允许创建 PR
```

#### 任务完成时的强制验证
```bash
# 当用户说"完成"、"done"时：
# ✅ 检查所有验收标准是否满足
# ✅ 检查 PR 状态和 CI/CD 状态
# ✅ 生成完成报告
# ✅ 提供后续步骤指导

# 如未完成会阻断并要求处理未完成项目
```

### 3. Hook 脚本说明

#### `validate-development-workflow.sh`
- **触发时机**: 每次工具使用后 (PostToolUse)
- **检查内容**: 代码质量、构建状态、测试覆盖率
- **阻断条件**: 任何检查失败时返回 exit code 2
- **输出**: 详细的检查结果和修复建议

#### `validate-pr-readiness.sh`  
- **触发时机**: PR 创建前 (PreToolUse)
- **检查内容**: 分支状态、工作目录、验收完成度
- **阻断条件**: 前置条件未满足时阻断 PR 创建
- **输出**: PR 创建准备度报告

#### `check-task-completion.sh`
- **触发时机**: 用户提到"完成"时 (UserPromptSubmit)  
- **检查内容**: 完整的任务完成度验证
- **阻断条件**: 任务未真正完成时强制要求处理
- **输出**: 任务完成报告和后续步骤指导

## 📊 验收标准检查清单

自动维护的检查状态文件: `.claude/workflow-checklist.json`

```json
{
  "checks": {
    "code_quality": "passed|failed",
    "build_test": "passed|failed", 
    "test_coverage": "passed|failed",
    "overall_status": "ready_for_pr|needs_fixes"
  }
}
```

## 🔧 自定义配置

### 修改覆盖率要求
编辑脚本中的覆盖率检查逻辑：
```bash
# 在 validate-development-workflow.sh 中
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    # 修改为其他阈值，如 85
```

### 添加自定义检查
在任何 hook 脚本中添加新的检查函数：
```bash
check_custom_requirement() {
    # 自定义检查逻辑
    if [ 自定义条件 ]; then
        log_info "自定义检查通过"
        return 0
    else
        log_error "自定义检查失败"
        return 1
    fi
}
```

### 禁用特定检查
在配置文件中注释或删除不需要的 hook：
```json
{
  "hooks": {
    "PostToolUse": [
      // 注释掉不需要的 hook
      // {
      //   "name": "development_workflow_validation",
      //   ...  
      // }
    ]
  }
}
```

## 🚨 故障排除

### Hook 不生效
1. 确认配置文件路径正确
2. 检查脚本执行权限: `chmod +x .claude/scripts/*.sh`
3. 查看 Claude Code 日志获取错误信息

### Hook 过度严格
1. 临时禁用: 重命名 hooks-config.json
2. 调整检查条件: 编辑对应的脚本文件
3. 添加例外处理: 在脚本中增加特殊情况处理

### 依赖缺失
脚本依赖以下工具:
- `jq`: JSON 处理 (`brew install jq`)
- `bc`: 数值计算 (通常已预装)
- `gh`: GitHub CLI
- Go 工具链: `go`, `gofmt`, `goimports`

## 💡 最佳实践

1. **定期更新检查标准**: 根据项目演进调整 hook 检查逻辑
2. **团队一致性**: 确保所有开发者使用相同的 hook 配置  
3. **渐进式启用**: 新项目可以从宽松检查开始，逐步严格化
4. **监控效果**: 定期检查 hook 阻断的问题类型，优化检查逻辑

---

通过这套 hooks 系统，Claude Code 将被强制执行完整的开发流程，确保没有任何验收标准被遗漏！