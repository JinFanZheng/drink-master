# Docker 部署指南

## 阿里云容器镜像服务部署

本项目已配置自动构建和推送Docker镜像到阿里云容器镜像服务（ACR）。

### 🔧 配置信息

- **镜像仓库**: `registry.cn-shenzhen.aliyuncs.com/lrmtc/drink-master`
- **支持平台**: `linux/amd64`
- **版本管理**: 基于 `VERSION` 文件进行语义化版本控制

### 🚀 快速开始

#### 1. 登录阿里云容器镜像服务

```bash
make docker-login
```

#### 2. 构建并推送镜像

```bash
# 方式一: 分步执行 (推荐用于本地开发和调试)
make docker-push

# 方式二: 直接构建推送 (推荐用于CI/CD，更高效)
make docker-build-and-push

# 方式三: 完整的发布流程（包含测试）
make release-current

# 方式四: 快速发布流程（推荐）
make release-current-fast
```

### 🔧 两种构建方式的区别

#### 分步构建 (`docker-push`)
- 先构建镜像到本地 Docker daemon (`--load`)
- 然后推送到阿里云镜像仓库
- 适合本地开发和调试
- 镜像会保留在本地

#### 直接构建推送 (`docker-build-and-push`)
- 直接构建并推送到阿里云镜像仓库 (`--push`)
- 不会保留本地镜像副本
- 更高效，适合CI/CD环境
- **推荐用于生产发布**

### 📋 可用命令

#### Docker 相关命令

```bash
# 构建本地开发镜像
make docker-build

# 构建生产环境镜像 (linux/amd64)
make docker-build-prod

# 构建并推送到阿里云 (分步)
make docker-push

# 直接构建并推送到阿里云 (一步完成，推荐)
make docker-build-and-push

# 登录阿里云容器镜像服务
make docker-login

# 运行本地容器
make docker-run

# 运行生产环境容器
make docker-run-prod
```

#### 版本管理命令

```bash
# 查看当前版本
make version

# 升级补丁版本 (1.0.0 -> 1.0.1)
make version-patch

# 升级次版本 (1.0.0 -> 1.1.0)
make version-minor

# 升级主版本 (1.0.0 -> 2.0.0)
make version-major

# 设置指定版本
make version-set NEW_VERSION=v1.2.3
```

#### 发布流程命令

```bash
# 补丁发布（测试 + 升级补丁版本 + 推送）
make release-patch

# 次版本发布（测试 + 升级次版本 + 推送）
make release-minor

# 主版本发布（测试 + 升级主版本 + 推送）
make release-major

# 发布当前版本（测试 + 推送当前版本）
make release-current
```

### 🔄 典型发布工作流

#### 日常Bug修复（补丁版本）

```bash
# 1. 开发和测试代码
git add .
git commit -m "fix: 修复订单状态更新问题"

# 2. 发布补丁版本
make release-patch
```

#### 新功能发布（次版本）

```bash
# 1. 开发和测试代码
git add .
git commit -m "feat: 添加用户积分功能"

# 2. 发布次版本
make release-minor
```

#### 重大更新（主版本）

```bash
# 1. 开发和测试代码
git add .
git commit -m "feat!: 重构API架构，不兼容旧版本"

# 2. 发布主版本
make release-major
```

### 🛠️ 部署到生产环境

#### 在服务器上拉取并运行

```bash
# 拉取最新镜像
docker pull registry.cn-shenzhen.aliyuncs.com/lrmtc/drink-master:latest

# 运行容器
docker run -d \
  --name drink-master \
  -p 8080:8080 \
  --env-file .env \
  --restart unless-stopped \
  registry.cn-shenzhen.aliyuncs.com/lrmtc/drink-master:latest
```

#### 使用Docker Compose

```yaml
version: '3.8'

services:
  drink-master:
    image: registry.cn-shenzhen.aliyuncs.com/lrmtc/drink-master:latest
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 📊 镜像信息查看

```bash
# 查看镜像详情
docker inspect registry.cn-shenzhen.aliyuncs.com/lrmtc/drink-master:latest

# 查看镜像层级
docker history registry.cn-shenzhen.aliyuncs.com/lrmtc/drink-master:latest
```

### ⚠️ 注意事项

1. **首次使用需要登录阿里云**: 运行 `make docker-login` 并输入阿里云凭证
2. **版本号格式**: 使用语义化版本 (vX.Y.Z)，如 v1.0.0
3. **平台兼容性**: 镜像专门为 linux/amd64 构建，适用于大多数云服务器
4. **自动化测试**: 发布命令会自动运行测试，确保代码质量
5. **健康检查**: 容器包含内置健康检查，监控 `/api/health` 端点

### 🔍 故障排除

#### 推送失败

```bash
# 检查登录状态
docker system info | grep Registry

# 重新登录
make docker-login
```

#### 构建失败

```bash
# 检查Docker版本（需要支持 buildx）
docker --version
docker buildx version

# 启用 buildx（如果需要）
docker buildx install
```

#### 版本冲突

```bash
# 检查当前版本
make version

# 手动设置版本
make version-set NEW_VERSION=v1.0.0
```