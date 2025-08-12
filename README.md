# 🏪 Smart Vending Machine Platform

> 基于Agent协作框架的智能售货机管理平台 - Go+Gin+MySQL

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-green.svg)](https://gin-gonic.com)
[![Database](https://img.shields.io/badge/Database-MySQL-orange.svg)](https://www.mysql.com)
[![WeChat Pay](https://img.shields.io/badge/Payment-WeChat-green.svg)](https://pay.weixin.qq.com)

## 🎯 项目简介

智能售货机平台是一个现代化的IoT设备管理和电商系统，为消费者提供便捷的饮品购买体验，为设备运营商提供高效的设备管理和盈利工具。

### 核心特性
- 🔐 微信登录和JWT认证系统
- 🏪 售货机设备管理和监控
- 🥤 商品管理和库存同步
- 📱 移动端用户购买体验
- 💰 微信支付集成和自动退款
- 📊 销售数据统计和运营分析
- 🔄 MQTT设备实时通信
- ⚡ 高性能API响应 (<500ms)

## 🏗️ 技术架构

### 技术栈
- **后端**: Go 1.21+ + Gin Framework
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT + 微信登录
- **支付**: 微信支付API
- **设备通信**: MQTT协议
- **配置**: 环境变量 + .env文件

### 项目结构
```
drink-master/
├── cmd/server/                  # 应用程序入口
├── internal/                    # 内部包，不对外暴露
│   ├── handlers/               # HTTP处理器 (Controller层)
│   │   ├── member.go           # 用户管理接口
│   │   ├── machine.go          # 售货机管理
│   │   ├── product.go          # 商品管理
│   │   ├── order.go            # 订单管理
│   │   └── payment.go          # 支付相关
│   ├── services/               # 业务逻辑层
│   ├── repositories/           # 数据访问层
│   ├── models/                 # 数据模型 (Entity层)
│   ├── contracts/              # API契约定义
│   └── middleware/             # Gin中间件
├── pkg/                        # 可复用的公共包
│   ├── wechat/                 # 微信SDK封装
│   └── mqtt/                   # MQTT客户端
├── migrations/                 # 数据库迁移脚本
├── docs/                       # 项目文档
│   └── PRD/                    # 产品需求文档
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
mysql -u root -p -e "CREATE DATABASE vending_machine_dev"

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

### 用户认证
```bash
# 检查用户登录状态
GET /api/Account/CheckLogin
Authorization: Bearer <token>

# 微信登录
POST /api/Account/WeChatLogin
{
  "appId": "wx1234567890",
  "code": "wx_js_code", 
  "avatarUrl": "https://avatar.url",
  "nickName": "用户昵称"
}

# 获取用户信息
GET /api/Account/GetUserInfo
Authorization: Bearer <token>

# 检查用户信息（通过code）
GET /api/Account/CheckUserInfo?code=wx_code&appId=wx_app_id
```

### 会员管理
```bash
# 更新会员信息
POST /api/Member/Update
Authorization: Bearer <token>
{
  "nickname": "新昵称",
  "avatar": "新头像URL"
}

# 添加加盟意向
POST /api/Member/AddFranchiseIntention
Authorization: Bearer <token>
{
  "contactName": "联系人",
  "contactPhone": "联系电话",
  "intendedLocation": "意向地点"
}
```

### 售货机管理
```bash
# 获取售货机分页列表
POST /api/Machine/GetPaging
Authorization: Bearer <token>
{
  "page": 1,
  "pageSize": 10,
  "keyword": "搜索关键词"
}

# 获取售货机列表
GET /api/Machine/GetList
Authorization: Bearer <token>

# 获取售货机详情
GET /api/Machine/Get?id=machine_id

# 检查设备是否存在
GET /api/Machine/CheckDeviceExist?deviceId=device_id

# 获取售货机商品列表
GET /api/Machine/GetProductList?machineId=machine_id

# 开关营业状态
GET /api/Machine/OpenOrClose?id=machine_id
Authorization: Bearer <token>
```

### 订单管理
```bash
# 获取我的订单列表
POST /api/Order/GetPaging
Authorization: Bearer <token>
{
  "page": 1,
  "pageSize": 10
}

# 获取订单详情
GET /api/Order/Get?id=order_id

# 创建订单
POST /api/Order/Create
Authorization: Bearer <token>
{
  "machineId": "售货机ID",
  "productId": "商品ID", 
  "hasCup": true,
  "quantity": 1
}

# 申请退款（机主权限）
POST /api/Order/Refund
Authorization: Bearer <token>
{
  "orderId": "订单ID",
  "refundReason": "退款原因"
}
```

### 支付管理
```bash
# 获取支付信息（发起支付）
GET /api/Payment/Get?orderId=order_id
Authorization: Bearer <token>

# 查询支付结果
GET /api/Payment/Query?orderId=order_id
Authorization: Bearer <token>
```

### 回调接口
```bash
# 微信支付结果回调
POST /api/Callback/PaymentResult
# 第三方支付平台调用，无需认证
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

### 会员 (Members)
- ID, Nickname, Avatar, WeChatOpenID
- Role (member/owner), MachineOwnerID, IsAdmin
- CreatedAt, UpdatedAt

### 设备运营商 (MachineOwners)  
- ID, Name, ContactPhone, ContactEmail
- ReceivingAccount, CreatedAt, UpdatedAt

### 售货机 (Machines)
- ID, MachineOwnerID, DeviceName, DeviceID
- Location, Status, IsBusinessOpen
- CreatedAt, UpdatedAt

### 商品 (Products)
- ID, Name, Description, ImageURL
- Category, CreatedAt, UpdatedAt

### 设备商品关联 (MachineProducts)
- ID, MachineID, ProductID
- Price, PriceWithoutCup, Stock, IsAvailable
- CreatedAt, UpdatedAt

### 订单 (Orders)
- ID, MemberID, MachineID, ProductID, OrderNo
- HasCup, TotalAmount, PayAmount
- PaymentStatus, MakeStatus, PaymentTime
- RefundAmount, RefundReason, CreatedAt

## 🔒 认证机制

系统使用微信登录 + JWT Token认证：

1. 用户通过微信小程序获取code，调用 `/api/Account/WeChatLogin` 接口登录
2. 系统验证微信code，创建或更新用户信息，返回JWT token
3. 后续请求在Header中携带token：`Authorization: Bearer <token>`
4. 受保护的路由通过JWT middleware验证token有效性
5. Token默认有效期24小时，支持机主和普通会员权限控制

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

**让我们通过标准化的协作流程，构建高质量的智能售货机平台！** 🏪

*最后更新：2025-08-11*