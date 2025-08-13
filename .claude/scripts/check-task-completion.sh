#!/bin/bash

# 任务完成检查 Hook
# 当用户说任务"完成"时，强制检查是否真正完成所有验收标准

set -e

PROJECT_ROOT=$(pwd)
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "main")
CHECKLIST_FILE=".claude/workflow-checklist.json"

log_info() {
    echo "📋 [Task Completion] $1" >&2
}

log_warning() {
    echo "⚠️ [Task Completion] $1" >&2
}

log_error() {
    echo "❌ [Task Completion] $1" >&2
}

# 检查任务是否真正完成
check_task_completion() {
    log_info "检查任务完成状态..."
    
    INCOMPLETE_ITEMS=()
    
    # 1. 检查是否在feature分支
    if [[ "$CURRENT_BRANCH" =~ ^feat/[0-9]+-.*$ ]]; then
        log_info "正在feature分支 $CURRENT_BRANCH 上工作"
        
        # 2. 检查代码质量
        if ! gofmt -l . | wc -l | grep -q "^0$"; then
            INCOMPLETE_ITEMS+=("代码格式化未完成")
        fi
        
        if ! goimports -d $(find . -name "*.go" -not -path "./vendor/*") 2>/dev/null | wc -l | grep -q "^0$"; then
            INCOMPLETE_ITEMS+=("import格式化未完成")
        fi
        
        # 3. 检查构建和测试
        if ! make build >/dev/null 2>&1; then
            INCOMPLETE_ITEMS+=("代码构建失败")
        fi
        
        if ! make test >/dev/null 2>&1; then
            INCOMPLETE_ITEMS+=("测试执行失败")
        fi
        
        if ! make lint >/dev/null 2>&1; then
            INCOMPLETE_ITEMS+=("Lint检查失败")
        fi
        
        # 4. 检查测试覆盖率
        if [ -f "coverage.out" ]; then
            COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
            if (( $(echo "$COVERAGE < 80" | bc -l) )); then
                INCOMPLETE_ITEMS+=("测试覆盖率不足: ${COVERAGE}% < 80%")
            fi
        else
            INCOMPLETE_ITEMS+=("缺少测试覆盖率报告")
        fi
        
        # 5. 检查PR状态
        PR_EXISTS=$(gh pr list --head "$CURRENT_BRANCH" --json number --jq length 2>/dev/null || echo "0")
        if [ "$PR_EXISTS" -eq 0 ]; then
            INCOMPLETE_ITEMS+=("尚未创建PR")
        else
            # 检查PR的CI/CD状态
            PR_STATUS=$(gh pr view --json statusCheckRollup --jq '.statusCheckRollup[0].conclusion // "PENDING"' 2>/dev/null || echo "UNKNOWN")
            if [ "$PR_STATUS" != "SUCCESS" ]; then
                INCOMPLETE_ITEMS+=("PR的CI/CD检查未通过，状态: $PR_STATUS")
            fi
        fi
        
        # 6. 检查是否有未提交的更改
        if ! git diff-index --quiet HEAD --; then
            INCOMPLETE_ITEMS+=("存在未提交的更改")
        fi
        
    else
        log_warning "不在feature分支上，跳过详细检查"
    fi
    
    # 输出检查结果
    if [ ${#INCOMPLETE_ITEMS[@]} -eq 0 ]; then
        log_info "✅ 任务完成检查通过！所有验收标准已满足。"
        
        # 提供下一步指导
        if [[ "$CURRENT_BRANCH" =~ ^feat/[0-9]+-.*$ ]]; then
            log_info "📋 后续步骤："
            log_info "1. 等待code review"
            log_info "2. PR合并后清理工作分支"
            log_info "3. 切换回主分支继续其他任务"
        fi
        
        return 0
    else
        log_error "⚠️ 任务尚未真正完成！以下项目需要处理："
        for item in "${INCOMPLETE_ITEMS[@]}"; do
            log_error "  - $item"
        done
        
        log_info "🔧 建议的修复步骤："
        log_info "1. 运行 'make lint && make test && make build' 进行质量检查"
        log_info "2. 运行 'go fmt ./...' 和 'find . -name \"*.go\" -not -path \"./vendor/*\" -exec goimports -w {} \\;' 修复格式"
        log_info "3. 添加测试用例提高覆盖率到80%以上"
        log_info "4. 提交所有更改并创建PR"
        log_info "5. 等待CI/CD检查全部通过"
        
        return 1
    fi
}

# 生成完成报告
generate_completion_report() {
    local report_file=".claude/task-completion-report.md"
    
    cat > "$report_file" << EOF
# 任务完成报告

**分支**: $CURRENT_BRANCH  
**检查时间**: $(date)

## 验收检查结果

EOF
    
    if check_task_completion >/dev/null 2>&1; then
        cat >> "$report_file" << EOF
✅ **状态**: 所有验收标准已满足

## 质量指标
- 代码格式: ✅ 通过
- 构建状态: ✅ 通过  
- 测试状态: ✅ 通过
- Lint检查: ✅ 通过
- 测试覆盖率: ✅ ≥80%
- PR状态: ✅ CI/CD通过

任务已准备好进行最终review和合并。
EOF
    else
        cat >> "$report_file" << EOF
❌ **状态**: 存在未完成项目

请按照上述建议完成所有验收标准后，重新检查。
EOF
    fi
    
    log_info "完成报告已生成: $report_file"
}

# 主函数
main() {
    log_info "开始任务完成验证..."
    
    if check_task_completion; then
        generate_completion_report
        log_info "🎉 恭喜！任务确实已完成所有验收标准。"
        exit 0
    else
        generate_completion_report
        log_error "⚠️ 请注意：任务尚未完全完成，请处理上述问题。"
        # 返回2阻断，强制Claude处理未完成的项目
        exit 2
    fi
}

main "$@"