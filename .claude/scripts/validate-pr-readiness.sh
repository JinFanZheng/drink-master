#!/bin/bash

# PR创建前验证脚本
# 确保所有质量检查完成才能创建PR

set -e

PROJECT_ROOT=$(pwd)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "main")
CHECKLIST_FILE=".claude/workflow-checklist.json"

log_info() {
    echo "🔍 [PR Validation] $1" >&2
}

log_error() {
    echo "❌ [PR Validation] $1" >&2
}

# 检查是否所有验收标准都通过
check_pr_readiness() {
    log_info "检查PR创建前置条件..."
    
    # 检查是否在正确的feature分支
    if [[ ! "$CURRENT_BRANCH" =~ ^feat/[0-9]+-.*$ ]]; then
        log_error "当前不在feature分支，无法创建PR"
        return 1
    fi
    
    # 检查工作目录是否干净
    if ! git diff-index --quiet HEAD --; then
        log_error "工作目录有未提交的更改，请先提交所有更改"
        return 1
    fi
    
    # 检查是否存在验收检查清单
    if [ ! -f "$CHECKLIST_FILE" ]; then
        log_error "未找到验收检查清单，请先完成开发验收检查"
        return 1
    fi
    
    # 检查验收状态
    if command -v jq >/dev/null 2>&1; then
        OVERALL_STATUS=$(jq -r '.checks.overall_status // "unknown"' "$CHECKLIST_FILE")
        if [ "$OVERALL_STATUS" != "ready_for_pr" ]; then
            log_error "验收检查未通过，当前状态: $OVERALL_STATUS"
            log_error "请先完成所有开发验收检查"
            return 1
        fi
    fi
    
    # 最终质量检查
    if ! make lint >/dev/null 2>&1; then
        log_error "Lint检查失败"
        return 1
    fi
    
    if ! make test >/dev/null 2>&1; then
        log_error "测试失败"
        return 1
    fi
    
    if ! make build >/dev/null 2>&1; then
        log_error "构建失败"
        return 1
    fi
    
    # 检查测试覆盖率
    if [ -f "coverage.out" ]; then
        COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            log_error "测试覆盖率不足: ${COVERAGE}% < 80%"
            return 1
        fi
    fi
    
    log_info "✅ PR创建前置条件检查通过"
    return 0
}

# 主函数
main() {
    if check_pr_readiness; then
        log_info "🚀 已准备好创建PR！"
        exit 0
    else
        log_error "🚫 PR创建被阻止，请完成上述检查后重试"
        exit 2  # 阻断操作
    fi
}

main "$@"