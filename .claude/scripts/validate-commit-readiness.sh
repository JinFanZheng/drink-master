#!/bin/bash

# 提交前验证脚本
# 确保代码质量和测试覆盖率达标才能提交

set -e

PROJECT_ROOT=$(pwd)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "main")

log_info() {
    echo "🔍 [Commit Validation] $1" >&2
}

log_error() {
    echo "❌ [Commit Validation] $1" >&2
}

# 检查提交准备度
check_commit_readiness() {
    log_info "检查提交前置条件..."
    
    # 1. 检查代码格式
    if ! gofmt -l . | wc -l | grep -q "^0$"; then
        log_error "代码格式不规范，请运行: go fmt ./..."
        return 1
    fi
    
    # 2. 检查goimports格式
    if ! goimports -d $(find . -name "*.go" -not -path "./vendor/*") 2>/dev/null | wc -l | grep -q "^0$"; then
        log_error "import格式不规范，请运行: find . -name \"*.go\" -not -path \"./vendor/*\" -exec goimports -w {} \\;"
        return 1
    fi
    
    # 3. 检查构建
    if ! make build >/dev/null 2>&1; then
        log_error "构建失败，请修复后重试"
        return 1
    fi
    
    # 4. 检查测试
    if ! make test >/dev/null 2>&1; then
        log_error "测试失败，请修复后重试"
        return 1
    fi
    
    # 5. 检查Lint
    if ! make lint >/dev/null 2>&1; then
        log_error "Lint检查失败，请修复后重试"
        return 1
    fi
    
    # 6. 检查测试覆盖率（仅在feature分支）
    if [[ "$CURRENT_BRANCH" =~ ^feat/[0-9]+-.*$ ]]; then
        if [ -f "coverage.out" ]; then
            COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
            if (( $(echo "$COVERAGE < 80" | bc -l) )); then
                log_error "测试覆盖率不足: ${COVERAGE}% < 80%"
                log_error "请添加更多测试用例后再提交"
                return 1
            else
                log_info "测试覆盖率检查通过: ${COVERAGE}%"
            fi
        else
            log_error "缺少测试覆盖率报告，请运行 'make test' 生成"
            return 1
        fi
    fi
    
    log_info "✅ 提交前置条件检查通过"
    return 0
}

# 主函数
main() {
    if check_commit_readiness; then
        log_info "🚀 代码已准备好提交！"
        exit 0
    else
        log_error "🚫 代码提交被阻止，请完成上述检查后重试"
        exit 2  # 阻断操作
    fi
}

main "$@"