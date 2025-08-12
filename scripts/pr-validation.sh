#!/bin/bash
# pr-validation.sh - PR合并前自动验证脚本
# 基于 docs/AGENT_PR_MERGE_GUIDE.md 实现

set -e

PR_NUMBER=$1
if [ -z "$PR_NUMBER" ]; then
  echo "Usage: $0 <pr-number>"
  exit 1
fi

echo "🔍 开始验证PR #$PR_NUMBER..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 输出函数
log_info() {
    echo -e "${BLUE}$1${NC}"
}

log_success() {
    echo -e "${GREEN}$1${NC}"
}

log_warning() {
    echo -e "${YELLOW}$1${NC}"
}

log_error() {
    echo -e "${RED}$1${NC}"
}

# 1. 基础检查
echo "1️⃣ 检查CI状态..."
if ! gh pr checks $PR_NUMBER | grep -q "✓"; then
  log_error "❌ CI检查未通过"
  exit 1
fi
log_success "✅ CI检查通过"

# 2. 检查可合并状态
echo "2️⃣ 检查合并状态..."
MERGEABLE=$(gh pr view $PR_NUMBER --json mergeable -q .mergeable)
if [ "$MERGEABLE" != "true" ]; then
  log_error "❌ PR有冲突，无法合并"
  exit 1
fi
log_success "✅ 无合并冲突"

# 3. 检查Issue链接
echo "3️⃣ 检查Issue关联..."
ISSUE_LINK=$(gh pr view $PR_NUMBER --json body -q .body | grep -E "(Fixes|Closes) #[0-9]+" || echo "")
if [ -z "$ISSUE_LINK" ]; then
  log_warning "⚠️ 未发现Issue链接"
  echo "建议在PR描述中添加 'Fixes #<issue-number>'"
else
  log_success "✅ Issue链接: $ISSUE_LINK"
fi

# 4. 代码质量检查
echo "4️⃣ 运行代码质量检查..."
log_info "执行 golangci-lint..."
if ! make lint > /dev/null 2>&1; then
  log_error "❌ 代码质量检查失败"
  echo "请运行 'make lint' 查看详细错误"
  exit 1
fi
log_success "✅ 代码质量检查通过"

# 4.1 检查测试覆盖率
echo "4.1️⃣ 检查测试覆盖率..."
log_info "运行测试并生成覆盖率报告..."
if ! make test > /dev/null 2>&1; then
  log_error "❌ 测试执行失败"
  exit 1
fi

if [ -f "coverage.out" ]; then
  COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
  if (( $(echo "$COVERAGE < 78" | bc -l) )); then
    log_error "❌ 测试覆盖率不足: ${COVERAGE}% (要求≥78%)"
    echo "请添加更多测试用例以提高覆盖率"
    exit 1
  fi
  log_success "✅ 测试覆盖率: ${COVERAGE}%"
else
  log_warning "⚠️ 未找到覆盖率报告文件"
fi

# 4.2 检查循环复杂度
echo "4.2️⃣ 检查循环复杂度..."
if golangci-lint run --disable-all --enable=gocyclo > /dev/null 2>&1; then
  log_success "✅ 循环复杂度检查通过"
else
  log_warning "⚠️ 发现循环复杂度过高的函数"
  echo "建议重构复杂函数为多个小函数"
  # 不阻塞合并，只是警告
fi

# 5. 构建检查
echo "5️⃣ 构建应用..."
if ! make build > /dev/null 2>&1; then
  log_error "❌ 应用构建失败"
  exit 1
fi
log_success "✅ 应用构建成功"

# 6. 功能验证（如果存在健康检查端点）
echo "6️⃣ 验证应用功能..."
if command -v make health-check &> /dev/null; then
  if ! make health-check > /dev/null 2>&1; then
    log_warning "⚠️ 健康检查失败，可能需要手动验证"
  else
    log_success "✅ 健康检查通过"
  fi
else
  log_info "跳过健康检查（命令不存在）"
fi

# 7. 检查变更类型和风险
echo "7️⃣ 分析变更风险..."
CHANGED_FILES=$(gh pr diff $PR_NUMBER --name-only)
PR_TITLE=$(gh pr view $PR_NUMBER --json title -q .title)
RISK_LEVEL="low"

log_info "变更文件:"
echo "$CHANGED_FILES" | head -10

# 检查关联Issue的里程碑
if [ ! -z "$ISSUE_LINK" ]; then
  ISSUE_NUM=$(echo "$ISSUE_LINK" | grep -oE '[0-9]+' | head -1)
  if [ ! -z "$ISSUE_NUM" ]; then
    MILESTONE=$(gh issue view $ISSUE_NUM --json milestone -q '.milestone.title // "无里程碑"' 2>/dev/null || echo "获取失败")
    log_info "📋 关联Issue #$ISSUE_NUM，里程碑: $MILESTONE"
  fi
fi

# 风险评估
if echo "$CHANGED_FILES" | grep -E "(internal/contracts|internal/handlers)" > /dev/null; then
  RISK_LEVEL="medium"
  log_warning "⚠️ 发现API或契约变更，风险等级: 中"
fi

if echo "$CHANGED_FILES" | grep -E "(migrations/|internal/models|security|auth)" > /dev/null; then
  RISK_LEVEL="high"
  log_error "🚨 发现高风险变更（数据库/安全），需要人工审核"
  exit 2
fi

# 基于PR标题的风险评估
if echo "$PR_TITLE" | grep -E "^(feat|fix|refactor):" > /dev/null; then
  if [ "$RISK_LEVEL" = "low" ]; then
    RISK_LEVEL="medium"
    log_warning "⚠️ 功能/修复/重构类型，风险等级: 中"
  fi
elif echo "$PR_TITLE" | grep -E "^(docs|style|test|chore):" > /dev/null; then
  log_info "📝 文档/样式/测试/维护类型变更"
fi

# 8. 最终检查和建议
echo "8️⃣ 最终验证结果..."
log_success "✅ PR #$PR_NUMBER 验证完成"
log_info "风险等级: $RISK_LEVEL"

case "$RISK_LEVEL" in
  "low")
    log_success "🟢 低风险，建议自动合并"
    echo "建议合并命令: gh pr merge $PR_NUMBER --squash --delete-branch"
    ;;
  "medium")
    log_warning "🟡 中风险，建议人工确认后合并"
    echo "建议合并命令: gh pr merge $PR_NUMBER --squash --delete-branch"
    ;;
  "high")
    log_error "🔴 高风险，必须人工审核"
    echo "请添加 'needs-human-review' 标签"
    exit 2
    ;;
esac

echo ""
log_success "验证完成! 所有检查通过 ✨"