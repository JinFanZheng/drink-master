#!/bin/bash

# 性能测试脚本
# 对应Issue #15的性能标准验证

set -e

echo "🚀 开始性能测试..."

# 确保服务在运行
echo "检查服务状态..."
if ! curl -f http://localhost:8080/api/health >/dev/null 2>&1; then
    echo "❌ 服务未启动，请先运行 'make dev'"
    exit 1
fi

echo "✅ 服务运行中"

# API响应时间测试
echo "🔍 测试API响应时间..."

# 登录接口 < 1000ms
echo "测试登录接口响应时间..."
curl -w "登录接口: %{time_total}s\n" -o /dev/null -s \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"code":"test123","nickName":"Test","avatarUrl":"http://example.com/avatar.jpg"}' \
    http://localhost:8080/api/Account/WeChatLogin

# 查询接口 < 500ms
echo "测试查询接口响应时间..."
curl -w "获取机器详情: %{time_total}s\n" -o /dev/null -s \
    http://localhost:8080/api/Machine/Get?id=test123

curl -w "检查设备存在: %{time_total}s\n" -o /dev/null -s \
    http://localhost:8080/api/Machine/CheckDeviceExist?deviceId=test123

curl -w "获取产品列表: %{time_total}s\n" -o /dev/null -s \
    http://localhost:8080/api/Machine/GetProductList?machineId=test123

curl -w "健康检查: %{time_total}s\n" -o /dev/null -s \
    http://localhost:8080/api/health

# 并发测试
echo "🔥 测试并发性能..."

# 使用ab工具测试并发
if command -v ab >/dev/null 2>&1; then
    echo "使用Apache Bench测试并发性能..."
    ab -n 1000 -c 10 http://localhost:8080/api/health
else
    echo "Apache Bench未安装，跳过并发测试"
fi

# 使用Go基准测试
echo "运行Go基准测试..."
go test -bench=BenchmarkAPIPerformance -benchmem ./internal/

echo "✅ 性能测试完成"