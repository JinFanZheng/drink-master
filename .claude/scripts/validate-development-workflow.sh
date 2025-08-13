#!/bin/bash

# Claude Code 开发流程验收检查 Hook
# 在每次工具执行后自动检查是否需要执行验收标准

set -e

PROJECT_ROOT=$(pwd)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "main")
CHECKLIST_FILE=".claude/workflow-checklist.json"

# 日志函数
log_info() {
    echo "🔍 [Workflow Hook] $1" >&2
}

log_warning() {
    echo "⚠️ [Workflow Hook] $1" >&2
}

log_error() {
    echo "❌ [Workflow Hook] $1" >&2
}

# 检查是否在feature分支
is_feature_branch() {
    [[ "$CURRENT_BRANCH" =~ ^feat/[0-9]+-.*$ ]]
}

# 检查测试覆盖率
check_test_coverage() {
    if [ -f "coverage.out" ]; then
        COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            log_error "测试覆盖率不足: ${COVERAGE}% < 80%"
            log_error "请添加更多测试用例后再提交代码"
            return 1
        else
            log_info "测试覆盖率检查通过: ${COVERAGE}%"
            return 0
        fi
    else
        log_warning "未找到覆盖率报告，请运行 'make test' 生成"
        return 1
    fi
}

# 检查代码质量
check_code_quality() {
    log_info "检查代码质量..."
    
    # 检查Go代码格式
    if ! gofmt -l . | wc -l | grep -q "^0$"; then
        log_error "代码格式不规范，请运行: go fmt ./..."
        return 1
    fi
    
    # 检查goimports格式
    if ! goimports -d $(find . -name "*.go" -not -path "./vendor/*") | wc -l | grep -q "^0$"; then
        log_error "import格式不规范，请运行: find . -name \"*.go\" -not -path \"./vendor/*\" -exec goimports -w {} \\;"
        return 1
    fi
    
    log_info "代码格式检查通过"
    return 0
}

# 检查构建状态
check_build_status() {
    log_info "检查构建状态..."
    
    if ! make build >/dev/null 2>&1; then
        log_error "构建失败，请运行 'make build' 检查错误"
        return 1
    fi
    
    if ! make test >/dev/null 2>&1; then
        log_error "测试失败，请运行 'make test' 检查错误"
        return 1
    fi
    
    if ! make lint >/dev/null 2>&1; then
        log_error "Lint检查失败，请运行 'make lint' 检查错误"
        return 1
    fi
    
    log_info "构建和测试检查通过"
    return 0
}

# 更新检查清单状态
update_checklist() {
    local check_name="$1"
    local status="$2"
    
    # 创建或更新检查清单
    if [ ! -f "$CHECKLIST_FILE" ]; then
        echo '{"checks": {}}' > "$CHECKLIST_FILE"
    fi
    
    # 使用jq更新状态（如果可用）
    if command -v jq >/dev/null 2>&1; then
        jq ".checks[\"$check_name\"] = \"$status\"" "$CHECKLIST_FILE" > "${CHECKLIST_FILE}.tmp" && mv "${CHECKLIST_FILE}.tmp" "$CHECKLIST_FILE"
    fi
}

# 主要验收检查逻辑
main() {
    log_info "开始开发流程验收检查..."
    
    # 只在feature分支上进行严格检查
    if ! is_feature_branch; then
        log_info "当前在$CURRENT_BRANCH分支，跳过严格验收检查"
        exit 0
    fi
    
    log_info "在feature分支 $CURRENT_BRANCH 上，执行完整验收检查"
    
    # 执行各项检查
    CHECKS_PASSED=true
    
    # 1. 代码质量检查
    if check_code_quality; then
        update_checklist "code_quality" "passed"
    else
        update_checklist "code_quality" "failed"
        CHECKS_PASSED=false
    fi
    
    # 2. 构建和测试检查
    if check_build_status; then
        update_checklist "build_test" "passed"
    else
        update_checklist "build_test" "failed"  
        CHECKS_PASSED=false
    fi
    
    # 3. 测试覆盖率检查
    if check_test_coverage; then
        update_checklist "test_coverage" "passed"
    else
        update_checklist "test_coverage" "failed"
        CHECKS_PASSED=false
    fi
    
    # 检查结果处理
    if [ "$CHECKS_PASSED" = true ]; then
        log_info "✅ 所有验收检查通过！"
        update_checklist "overall_status" "ready_for_pr"
        exit 0
    else
        log_error "❌ 验收检查失败！请修复上述问题后再继续。"
        update_checklist "overall_status" "needs_fixes"
        
        # 提供修复建议
        log_info "🔧 修复建议："
        log_info "1. 运行 'make lint && make test && make build' 进行本地验证"
        log_info "2. 运行 'go tool cover -func=coverage.out | tail -1' 检查测试覆盖率"
        log_info "3. 运行 'go fmt ./...' 修复代码格式"
        log_info "4. 运行 'find . -name \"*.go\" -not -path \"./vendor/*\" -exec goimports -w {} \\;' 修复import格式"
        
        # 返回2表示阻断后续操作
        exit 2
    fi
}

# 运行主函数
main "$@"