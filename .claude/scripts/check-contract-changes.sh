#!/bin/bash

# 契约变更检查 Hook
# 当修改internal/contracts时提醒注意事项

set -e

# 读取hook参数（JSON格式）
if [ -t 0 ]; then
    HOOK_INPUT=""
else
    HOOK_INPUT=$(cat)
fi

log_info() {
    echo "📋 [Contract Check] $1" >&2
}

log_warning() {
    echo "⚠️ [Contract Check] $1" >&2
}

# 主函数
main() {
    log_warning "检测到契约变更！"
    log_warning "请注意以下事项："
    log_info "1. 确保所有相关的handler和service已同步更新"
    log_info "2. 更新API文档说明破坏性变更"
    log_info "3. 考虑版本兼容性"
    log_info "4. 确保测试用例覆盖新的契约"
    
    # 不阻断，只是提醒
    exit 0
}

main "$@"