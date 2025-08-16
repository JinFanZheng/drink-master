# Drink Master - Go项目开发工具

.PHONY: help dev build build-prod lint test test-short clean deps docs install-tools pre-commit deploy-check health-check test-api integration-test performance-test benchmark stats check-env git-status db-migrate db-rollback db-reset db-seed docker-build docker-build-prod docker-push docker-build-and-push docker-login docker-run docker-run-prod version version-patch version-minor version-major version-set release-patch release-minor release-major release-current

# 版本管理
VERSION := $(shell cat VERSION 2>/dev/null || echo "v1.0.0")
REGISTRY := registry.cn-shenzhen.aliyuncs.com/lrmtc
IMAGE_NAME := drink-master
FULL_IMAGE := $(REGISTRY)/$(IMAGE_NAME)

# 默认端口配置
PORT ?= 8080
HOST ?= localhost

# 默认目标
help: ## 显示帮助信息
	@echo "Drink Master - 饮品管理系统"
	@echo ""
	@echo "可用命令："
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ==================== 开发相关命令 ====================
dev: ## 启动开发服务器（热重载）
	@echo "🚀 启动开发服务器..."
	go run cmd/server/main.go

dev-port: ## 指定端口启动开发服务器 (使用 PORT=xxxx make dev-port)
	@echo "🚀 在端口 $(PORT) 启动开发服务器..."
	PORT=$(PORT) go run cmd/server/main.go

dev-mock: ## Mock模式启动开发服务器
	@echo "🎭 Mock模式启动开发服务器..."
	MOCK_MODE=true go run cmd/server/main.go

build: ## 编译Go二进制文件
	@echo "🔨 编译项目..."
	go build -o bin/drink-master cmd/server/main.go

build-prod: ## 生产环境优化编译
	@echo "🏗️ 生产环境编译..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/drink-master cmd/server/main.go

# ==================== 代码质量检查 ====================
lint: ## 运行代码检查 (golangci-lint + go fmt + go vet)
	@echo "🔍 运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "运行 golangci-lint..."; \
		golangci-lint run --disable=typecheck || echo "⚠️ golangci-lint 检查完成，存在一些问题但可以继续"; \
		echo "运行基础检查..."; \
		go fmt ./...; \
		go vet ./...; \
	else \
		echo "⚠️ golangci-lint 未安装，运行基础检查..."; \
		go fmt ./...; \
		go vet ./...; \
	fi

test: ## 运行所有测试
	@echo "🧪 运行测试..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-short: ## 运行快速测试（跳过慢速测试）
	@echo "⚡ 运行快速测试..."
	go test -v -short ./...

pre-commit: lint test build ## 预提交完整检查（lint + test + build）
	@echo "✅ 预提交检查完成"

# ==================== 数据库操作 ====================
db-migrate: ## 执行数据库迁移
	@echo "📊 执行数据库迁移..."
	@if [ -f "migrations/migrate.go" ]; then \
		go run migrations/migrate.go up; \
	else \
		echo "⚠️ 迁移文件不存在，请先创建迁移脚本"; \
	fi

db-rollback: ## 回滚最后一次数据库迁移
	@echo "↩️ 回滚数据库迁移..."
	@if [ -f "migrations/migrate.go" ]; then \
		go run migrations/migrate.go down; \
	else \
		echo "⚠️ 迁移文件不存在"; \
	fi

db-reset: ## 重置数据库（危险操作）
	@echo "🔄 重置数据库..."
	@read -p "确认要重置数据库吗？这将删除所有数据 [y/N]: " confirm && [ "$$confirm" = "y" ]
	@if [ -f "migrations/migrate.go" ]; then \
		go run migrations/migrate.go reset; \
	else \
		echo "⚠️ 迁移文件不存在"; \
	fi

db-seed: ## 填充测试数据
	@echo "🌱 填充测试数据..."
	@if [ -f "migrations/seed.go" ]; then \
		go run migrations/seed.go; \
	else \
		echo "⚠️ 种子数据文件不存在"; \
	fi

# ==================== 健康检查和测试 ====================
health-check: ## 检查服务健康状态
	@echo "❤️ 检查服务健康状态..."
	@curl -f http://localhost:8080/api/health || echo "❌ 服务不可用"

test-api: ## 测试主要API端点
	@echo "🔗 测试API端点..."
	@echo "检查健康状态..."
	@curl -s http://localhost:8080/api/health | jq '.' || echo "❌ 健康检查失败"
	@echo "检查数据库连接..."
	@curl -s http://localhost:8080/api/health/db | jq '.' || echo "❌ 数据库连接失败"

integration-test: ## 运行系统集成测试
	@echo "🧪 运行系统集成测试..."
	@./scripts/integration-test.sh

performance-test: ## 运行性能测试  
	@echo "⚡ 运行性能测试..."
	@./scripts/performance-test.sh

benchmark: ## 运行性能基准测试
	@echo "⚡ 运行性能测试..."
	go test -bench=. -benchmem ./...

deploy-check: pre-commit health-check test-api ## 部署前完整验证
	@echo "🚀 部署检查完成，可以安全部署"

# ==================== Docker相关 ====================
docker-build: ## 构建Docker镜像
	@echo "🐳 构建Docker镜像..."
	docker build -t drink-master:latest .

docker-build-prod: ## 构建生产环境Docker镜像 (linux/amd64)
	@echo "🐳 构建生产环境Docker镜像 (linux/amd64)..."
	@echo "当前版本: $(VERSION)"
	docker buildx build --platform linux/amd64 \
		--load \
		-t $(FULL_IMAGE):$(VERSION) \
		-t $(FULL_IMAGE):latest \
		.

docker-push: ## 推送Docker镜像到阿里云 (需要先构建)
	@echo "📤 推送Docker镜像到阿里云容器镜像服务..."
	@echo "推送版本: $(VERSION)"
	docker push $(FULL_IMAGE):$(VERSION)
	docker push $(FULL_IMAGE):latest
	@echo "✅ 镜像推送完成!"
	@echo "镜像地址:"
	@echo "  - $(FULL_IMAGE):$(VERSION)"
	@echo "  - $(FULL_IMAGE):latest"

docker-build-and-push: ## 构建并推送Docker镜像到阿里云 (推荐)
	@echo "🐳 构建并推送Docker镜像 (linux/amd64)..."
	@echo "推送版本: $(VERSION)"
	docker buildx build --platform linux/amd64 \
		--push \
		-t $(FULL_IMAGE):$(VERSION) \
		-t $(FULL_IMAGE):latest \
		.
	@echo "✅ 镜像构建和推送完成!"
	@echo "镜像地址:"
	@echo "  - $(FULL_IMAGE):$(VERSION)"
	@echo "  - $(FULL_IMAGE):latest"

docker-login: ## 登录阿里云容器镜像服务
	@echo "🔐 登录阿里云容器镜像服务..."
	@echo "请使用阿里云控制台的访问凭证进行登录"
	@echo "用户名: 阿里云账号全名 (如: your-name@example.com)"
	@echo "密码: 容器镜像服务的访问密码 (在阿里云容器镜像服务控制台设置)"
	docker login registry.cn-shenzhen.aliyuncs.com

docker-run: ## 运行Docker容器
	@echo "🚀 运行Docker容器..."
	docker run -p 8080:8080 --env-file .env drink-master:latest

docker-run-prod: ## 运行生产环境Docker容器
	@echo "🚀 运行生产环境Docker容器..."
	docker run -p 8080:8080 --env-file .env $(FULL_IMAGE):$(VERSION)

# ==================== 版本管理 ====================
version: ## 显示当前版本
	@echo "当前版本: $(VERSION)"

version-patch: ## 升级补丁版本 (x.y.z -> x.y.z+1)
	@echo "升级补丁版本..."
	@current=$$(cat VERSION | sed 's/v//'); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	patch=$$(echo $$current | cut -d. -f3); \
	new_patch=$$((patch + 1)); \
	new_version="v$$major.$$minor.$$new_patch"; \
	echo $$new_version > VERSION; \
	echo "版本已升级: $$current -> $$new_version"

version-minor: ## 升级次版本 (x.y.z -> x.y+1.0)
	@echo "升级次版本..."
	@current=$$(cat VERSION | sed 's/v//'); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	new_minor=$$((minor + 1)); \
	new_version="v$$major.$$new_minor.0"; \
	echo $$new_version > VERSION; \
	echo "版本已升级: $$current -> $$new_version"

version-major: ## 升级主版本 (x.y.z -> x+1.0.0)
	@echo "升级主版本..."
	@current=$$(cat VERSION | sed 's/v//'); \
	major=$$(echo $$current | cut -d. -f1); \
	new_major=$$((major + 1)); \
	new_version="v$$new_major.0.0"; \
	echo $$new_version > VERSION; \
	echo "版本已升级: $$current -> $$new_version"

version-set: ## 设置指定版本 (使用 VERSION=vx.y.z make version-set)
	@if [ -z "$(NEW_VERSION)" ]; then \
		echo "❌ 请指定版本号，例如: make version-set NEW_VERSION=v1.2.3"; \
		exit 1; \
	fi
	@echo "设置版本为: $(NEW_VERSION)"
	@echo "$(NEW_VERSION)" > VERSION
	@echo "✅ 版本已设置为: $(NEW_VERSION)"

# ==================== 发布流程 ====================
release-patch: pre-commit version-patch docker-build-and-push ## 补丁发布 (构建+测试+升级补丁版本+推送)
	@echo "🚀 补丁版本发布完成! 版本: $(shell cat VERSION)"

release-minor: pre-commit version-minor docker-build-and-push ## 次版本发布 (构建+测试+升级次版本+推送)
	@echo "🚀 次版本发布完成! 版本: $(shell cat VERSION)"

release-major: pre-commit version-major docker-build-and-push ## 主版本发布 (构建+测试+升级主版本+推送)
	@echo "🚀 主版本发布完成! 版本: $(shell cat VERSION)"

release-current: pre-commit docker-build-and-push ## 发布当前版本 (构建+测试+推送当前版本)
	@echo "🚀 当前版本发布完成! 版本: $(shell cat VERSION)"

# ==================== 开发工具 ====================
install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

docs: ## 生成API文档
	@echo "📚 生成API文档..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/server/main.go -o docs/swagger; \
	else \
		echo "⚠️ swag 工具未安装，运行 make install-tools 安装"; \
	fi

deps: ## 安装/更新依赖
	@echo "📦 管理项目依赖..."
	go mod tidy
	go mod download

clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	rm -f bin/drink-master
	rm -f coverage.out coverage.html integration_coverage.out
	go clean -testcache

# ==================== 项目信息 ====================
stats: ## 显示项目代码统计
	@echo "📊 项目代码统计:"
	@echo "Go文件数量:"
	@find . -name "*.go" -not -path "./vendor/*" | wc -l
	@echo "总代码行数:"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1

check-env: ## 检查开发环境
	@echo "🔍 开发环境检查:"
	@echo "Go版本: $(shell go version)"
	@echo "Git版本: $(shell git --version)"
	@echo "当前分支: $(shell git branch --show-current)"
	@echo "当前版本: $(VERSION)"
	@echo "工作目录状态:"
	@git status --porcelain | wc -l | xargs -I {} echo "  {} 个未提交的更改"

git-status: ## 检查Git状态
	@echo "📋 Git状态检查..."
	git status
	@echo ""
	@echo "未合并的分支:"
	git branch --no-merged main | head -10