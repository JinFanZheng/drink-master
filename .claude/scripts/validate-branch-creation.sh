#!/bin/bash

# 分支创建验证 Hook
# 确保从正确的基础分支创建feature分支

set -e

# 读取hook参数（JSON格式）
if [ -t 0 ]; then
    HOOK_INPUT=""
else
    HOOK_INPUT=$(cat)
fi

PROJECT_ROOT=$(pwd)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "main")

log_info() {
    echo "🔍 [Branch Validation] $1" >&2
}

log_error() {
    echo "❌ [Branch Validation] $1" >&2
}

# 主验证函数
validate_branch_creation() {
    log_info "验证分支创建前置条件..."
    
    # 1. 必须从main分支创建新分支
    if [ "$CURRENT_BRANCH" != "main" ]; then
        log_error "错误：必须从main分支创建feature分支！"
        log_error "当前分支：$CURRENT_BRANCH"
        log_error "请执行：git checkout main && git pull origin main"
        exit 2
    fi
    
    # 2. 确保main分支是最新的
    git fetch origin main --quiet
    LOCAL=$(git rev-parse main)
    REMOTE=$(git rev-parse origin/main)
    
    if [ "$LOCAL" != "$REMOTE" ]; then
        log_error "main分支不是最新的！"
        log_error "请执行：git pull origin main"
        exit 2
    fi
    
    # 3. 工作目录必须干净
    if ! git diff-index --quiet HEAD --; then
        log_error "工作目录有未提交的更改！"
        log_error "请先提交或stash当前更改"
        git status --short
        exit 2
    fi
    
    log_info "✅ 分支创建前置条件满足"
    log_info "- 基于最新main分支"
    log_info "- 工作目录干净"
    log_info "- 可以安全创建feature分支"
    
    return 0
}

# 主函数
main() {
    validate_branch_creation
    exit 0
}

main "$@"