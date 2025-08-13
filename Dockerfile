# 多阶段构建Dockerfile

# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装ca-certificates和curl用于健康检查和集成测试
RUN apk --no-cache add ca-certificates curl wget

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件和脚本
COPY --from=builder /app/.env.example .env
COPY --from=builder /app/scripts/ ./scripts/

# 更改文件所有者
RUN chown -R appuser:appgroup /app && \
    chmod +x ./scripts/*.sh

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查 - 验证API响应时间 < 500ms
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f --max-time 0.5 http://localhost:8080/api/health || exit 1

# 环境变量
ENV GIN_MODE=release
ENV LOG_LEVEL=info

# 启动应用
CMD ["./main"]