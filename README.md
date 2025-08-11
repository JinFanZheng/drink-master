# 🍹 Drink Master

> 基于Agent协作框架的Go+Gin+MySQL饮品管理系统

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-green.svg)](https://gin-gonic.com)
[![Database](https://img.shields.io/badge/Database-MySQL-orange.svg)](https://www.mysql.com)

## 🎯 项目简介

Drink Master是一个现代化的饮品管理系统，采用契约优先开发模式和Agent协作框架，提供完整的饮品记录、统计分析和用户管理功能。

### 核心特性
- 🔐 JWT用户认证系统
- 🍺 完整的饮品CRUD操作
- 📊 消费统计和趋势分析
- 🏷️ 饮品分类管理
- ⚡ 高性能API响应 (<500ms)
- 📝 自动生成API文档
- 🔄 支持Mock模式开发

## 🏗️ 技术架构

### 技术栈
- **后端**: Go 1.21+ + Gin Framework
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT (golang-jwt/jwt)
- **配置**: 环境变量 + .env文件

### 项目结构
```
drink-master/
├── cmd/server/                  # 应用程序入口
├── internal/                    # 内部包，不对外暴露
│   ├── handlers/               # HTTP处理器 (Controller层)
│   ├── services/               # 业务逻辑层 (Service层)
│   ├── repositories/           # 数据访问层 (Repository层)
│   ├── models/                 # 数据模型 (Entity层)
│   ├── contracts/              # API契约定义
│   └── middleware/             # Gin中间件
├── pkg/                        # 可复用的公共包
├── migrations/                 # 数据库迁移脚本
├── docs/                       # 项目文档
├── Makefile                    # 开发工具命令
└── CLAUDE.md                   # Agent协作指南
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Git

### 1. 克隆项目
```bash
git clone https://github.com/ddteam/drink-master.git
cd drink-master
```

### 2. 环境配置
```bash
# 复制环境配置文件
cp .env.example .env

# 编辑配置文件，设置数据库连接等
vi .env
```

### 3. 安装依赖
```bash
# 安装Go依赖
go mod tidy

# 安装开发工具（可选）
make install-tools
```

### 4. 数据库准备
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE drink_master_dev"

# 运行数据库迁移
make db-migrate

# 填充测试数据（可选）
make db-seed
```

### 5. 启动服务
```bash
# 开发模式启动
make dev

# 或者使用标准go命令
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

### 6. 验证安装
```bash
# 健康检查
curl http://localhost:8080/api/health

# API测试
make test-api
```

## 🔧 开发命令

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make dev` | 启动开发服务器（热重载） |
| `make build` | 编译Go二进制文件 |
| `make lint` | 代码质量检查 |
| `make test` | 运行所有测试 |
| `make db-migrate` | 执行数据库迁移 |
| `make health-check` | 检查服务健康状态 |
| `make pre-commit` | 预提交完整检查 |

### 完整命令列表
```bash
# 开发相关
make dev              # 启动开发服务器
make dev-mock         # Mock模式启动
make build            # 编译项目
make clean            # 清理构建文件

# 代码质量
make lint             # 代码检查
make test             # 运行测试
make pre-commit       # 预提交检查

# 数据库操作
make db-migrate       # 执行迁移
make db-rollback      # 回滚迁移
make db-reset         # 重置数据库
make db-seed          # 填充测试数据

# 健康检查
make health-check     # 服务健康检查
make test-api         # API功能测试
make deploy-check     # 部署前完整验证
```

## 📡 API接口

### 认证相关
```bash
# 用户注册
POST /api/auth/register
{
  "username": "user123",
  "email": "user@example.com", 
  "password": "password123"
}

# 用户登录
POST /api/auth/login
{
  "username": "user123",
  "password": "password123"
}
```

### 饮品管理
```bash
# 获取饮品列表
GET /api/drinks?category=coffee&limit=10&offset=0

# 创建饮品记录
POST /api/drinks
{
  "name": "拿铁咖啡",
  "category": "coffee",
  "price": 25.5,
  "description": "香浓拿铁"
}

# 获取单个饮品
GET /api/drinks/:id

# 更新饮品信息
PUT /api/drinks/:id

# 删除饮品
DELETE /api/drinks/:id
```

### 消费统计
```bash
# 消费统计
GET /api/stats/consumption?period=week

# 热门饮品
GET /api/stats/popular?limit=10

# 消费趋势
GET /api/stats/trends?period=month
```

### 系统相关
```bash
# 健康检查
GET /api/health

# 数据库健康检查  
GET /api/health/db
```

## 🧪 测试

### 运行测试
```bash
# 运行所有测试
make test

# 快速测试（跳过慢速测试）
make test-short

# 性能基准测试
make benchmark
```

### Mock模式测试
```bash
# 启用Mock模式
MOCK_MODE=true make dev

# 或者使用make命令
make dev-mock
```

## 📊 数据模型

### 用户 (User)
- ID, Username, Email, Password
- CreatedAt, UpdatedAt

### 饮品 (Drink) 
- ID, Name, Category, Price, Description
- UserID (外键), CreatedAt, UpdatedAt

### 饮品分类 (DrinkCategory)
- ID, Name, Description
- CreatedAt, UpdatedAt

### 消费记录 (ConsumptionLog)
- ID, DrinkID, UserID, ConsumedAt
- Quantity, Notes, CreatedAt

## 🔒 认证机制

系统使用JWT (JSON Web Token) 进行用户认证：

1. 用户通过 `/api/auth/login` 获取JWT token
2. 后续请求在Header中携带token：`Authorization: Bearer <token>`
3. 受保护的路由会验证token有效性
4. Token默认有效期24小时

## 🚀 部署

### 生产环境编译
```bash
make build-prod
```

### Docker部署
```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### 环境变量配置
生产环境必须设置的环境变量：
```bash
GIN_MODE=release
DB_HOST=production-mysql-host
DB_PASSWORD=secure-password
JWT_SECRET=production-jwt-secret
```

## 📚 Agent协作开发

本项目采用标准化的Agent协作框架，详细开发规范请参考：

- **[CLAUDE.md](CLAUDE.md)** - Agent协作总指南
- **[docs/README.md](docs/README.md)** - 完整文档导航
- **[docs/AGENT_ONBOARDING.md](docs/AGENT_ONBOARDING.md)** - Dev Agent开发流程
- **[docs/ROLES_COLLABORATION.md](docs/ROLES_COLLABORATION.md)** - 角色协作框架

### 核心开发流程
```bash
# 标准开发流程（严格执行）
git checkout main && git pull origin main
git status  # 确认工作目录干净
gh issue view <issue-id>
git checkout -b feat/<issue-id>-<name>
make lint && make test && make build

# 开发完成后
git commit -m "feat: implement feature"
gh pr create --title "feat: feature" --body "Fixes #<issue-id>"
```

## 📈 质量标准

### 代码质量要求
- ✅ Lint检查通过: `make lint`
- ✅ 所有测试通过: `make test`
- ✅ 构建成功: `make build`
- ✅ 测试覆盖率 > 80%

### API性能要求
- ✅ 响应时间 < 500ms
- ✅ 健康检查可用: `/api/health`
- ✅ 数据库连接正常: `/api/health/db`
- ✅ 并发支持: 1000+ req/s

### 验收标准 (Definition of Done)
- [ ] 功能完整实现且符合需求
- [ ] 单元测试覆盖率 > 80%
- [ ] API响应时间 < 500ms
- [ ] 数据库事务一致性保证
- [ ] CI/CD流水线全部通过
- [ ] 相关文档同步更新

## 🤝 贡献指南

1. Fork本仓库
2. 创建功能分支: `git checkout -b feat/amazing-feature`
3. 遵循开发规范: 参考 `docs/AGENT_ONBOARDING.md`
4. 提交代码: `git commit -m 'feat: add amazing feature'`
5. 推送分支: `git push origin feat/amazing-feature`
6. 创建Pull Request

### 提交规范
遵循 Conventional Commits 格式：
- `feat:` 新功能
- `fix:` Bug修复  
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `test:` 测试相关
- `chore:` 工具配置更新

## 📜 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🔗 相关链接

- [项目文档](docs/README.md)
- [API文档](docs/swagger/) (运行 `make docs` 生成)
- [开发指南](docs/AGENT_ONBOARDING.md)
- [协作框架](docs/ROLES_COLLABORATION.md)

---

**让我们通过标准化的协作流程，构建高质量的饮品管理系统！** 🍹

*最后更新：2025-08-11*