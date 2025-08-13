# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 🎯 Intent Recognition & Guide Loading

### 讨论需求 (Product Discussion)
**User says**: "我想讨论一下用户登录功能" / "Let's discuss the user authentication feature"
**Claude should**: 
1. Read `docs/WORKFLOWS/REQUIREMENT.md` for requirement analysis workflow
2. Analyze requirements and create PRD
3. Create PRD in `docs/PRD/` following template
4. Define acceptance criteria and success metrics

**User says**: "这个功能的用户场景是什么" / "What are the user scenarios for this feature?"
**Claude should**: 
1. Read existing PRDs in `docs/PRD/` for format reference
2. Perform user research and identify use cases
3. Document in standard PRD format

### 处理任务 (Task Processing)
**User says**: "处理 #54" / "Work on issue #54" / "开始做 #54"
**Claude should**: 
1. Read `docs/WORKFLOWS/DEVELOPMENT.md` for complete Dev workflow
2. View issue details: `gh issue view 54`
3. Check dependencies using `docs/WORKFLOWS/TASK_PARALLEL.md` if complex task
4. Follow branch creation workflow from guide
5. Use TodoWrite to plan implementation
6. Start development following quality standards

**User says**: "继续之前的任务" / "Continue the previous task"
**Claude should**: 
1. Check current branch: `git branch --show-current`
2. If on feature branch, continue from `docs/WORKFLOWS/DEVELOPMENT.md` section 5
3. Review uncommitted changes and resume

### 审查PR (PR Review)
**User says**: "审查 PR #123" / "Review PR #123" / "看一下这个PR"
**Claude should**:
1. Read `docs/WORKFLOWS/PR_MERGE.md` for complete review checklist
2. Execute validation steps from section "2. 合并前检查清单"
3. Run automated validation script if available
4. Classify risk level (Low/Medium/High)
5. Provide review feedback or merge decision

**User says**: "合并这个PR" / "Merge this PR"
**Claude should**: 
1. Read `docs/WORKFLOWS/PR_MERGE.md` section "3. 风险分类与合并策略"
2. Validate all required checks pass
3. Choose appropriate merge strategy based on risk
4. Execute merge command

### 任务管理 (Task Management)
**User says**: "创建一个Epic" / "Create an Epic" / "拆解任务"
**Claude should**:
1. Read `docs/WORKFLOWS/REQUIREMENT.md` section 2 for task breakdown
2. Read `docs/WORKFLOWS/TASK_PARALLEL.md` for parallel analysis
3. Create Epic with proper dependency marking

**User says**: "分析任务依赖" / "Analyze task dependencies" / "哪些可以并行"
**Claude should**:
1. Read `docs/WORKFLOWS/TASK_PARALLEL.md` completely
2. Check file conflicts
3. Identify parallel execution opportunities
4. Create task groups

### 紧急处理 (Emergency)
**User says**: "紧急修复" / "生产bug" / "hotfix"
**Claude should**:
1. Read `docs/EMERGENCY/HOTFIX.md`
2. Create hotfix branch from main
3. Apply minimal fix
4. Create urgent PR

**User says**: "需要回滚" / "rollback" / "恢复之前版本"
**Claude should**:
1. Read `docs/EMERGENCY/ROLLBACK.md`
2. Identify rollback scope
3. Execute appropriate rollback

### 查看项目状态 (Project Status)
**User says**: "现在有哪些待处理的任务" / "What tasks are pending?"
**Claude should**: 
1. Run `gh issue list --state open --label backend,api`
2. Provide summary of open issues

**User says**: "有哪些被阻塞的任务" / "What tasks are blocked?"
**Claude should**: 
1. Run `gh issue list --label blocked --state open`
2. Read each blocked issue to understand blockers
3. Suggest resolution strategies

## 📖 Keyword-to-Guide Mapping

When user mentions these keywords, automatically read the corresponding guide:

| Keywords | Guide to Read | Purpose |
|----------|--------------|---------|
| 处理任务, work on issue, 开发, development | `docs/WORKFLOWS/DEVELOPMENT.md` | Dev workflow |
| 审查PR, review PR, 合并, merge | `docs/WORKFLOWS/PR_MERGE.md` | PR validation |
| 需求, requirement, PRD, 产品设计 | `docs/WORKFLOWS/REQUIREMENT.md` | Requirement analysis |
| 拆解, Epic, 任务管理 | `docs/WORKFLOWS/REQUIREMENT.md` | Task breakdown |
| 并行, parallel, 依赖, dependency | `docs/WORKFLOWS/TASK_PARALLEL.md` | Parallel analysis |
| 紧急, urgent, hotfix, 生产bug | `docs/EMERGENCY/HOTFIX.md` | Emergency fix |
| 回滚, rollback, revert | `docs/EMERGENCY/ROLLBACK.md` | Rollback guide |
| API设计, 接口 | `docs/PATTERNS/API_DESIGN.md` | API patterns |
| 错误处理, error | `docs/PATTERNS/ERROR_HANDLING.md` | Error patterns |

## 🚀 Quick Start Guide

**Get started in 1 minute with the most common operations:**

### 📚 Required Reading
**Quick Start**: `docs/QUICK_START.md` - Get started in 1 minute
**Core Workflows**: 
- Development → `docs/WORKFLOWS/DEVELOPMENT.md`
- Requirements → `docs/WORKFLOWS/REQUIREMENT.md`
- PR Review → `docs/WORKFLOWS/PR_MERGE.md`
- Parallel Tasks → `docs/WORKFLOWS/TASK_PARALLEL.md`

### 1. Development Workflow (MANDATORY)

#### Starting a New Task
```bash
# When user says: "处理 #54" or "Work on issue #54"
git checkout main && git pull origin main           # Ensure latest code
git status                                          # Confirm clean working directory
gh issue view <issue-id>                           # Review task requirements
gh issue edit <issue-id> --add-label "in-progress" # Mark as in-progress
git checkout -b feat/<issue-id>-<short-name>       # Create feature branch
make lint && make test && make build               # Basic quality checks
```

#### Resuming a Task
```bash
# When user says: "继续任务" or "Continue the task"
git branch --show-current                          # Check current branch
git status                                         # Check uncommitted changes
gh issue view <issue-id>                          # Review requirements again
# Continue from where left off
```

### 2. Task Planning (MANDATORY)
- **MUST use TodoWrite tool** to create detailed task plan
- **MUST** break down complex tasks into specific steps
- **MUST** update task progress in real-time (pending → in_progress → completed)

### 3. Development Implementation Standards
- **Contract-First**: If modifying `internal/contracts/*.go`, must open PR and document breaking changes
- **Type Safety**: Strict Go typing, avoid `interface{}` when possible
- **Quality Gates**: Run `make lint && make test` after major changes

### 4. Pre-Commit Validation (MANDATORY)
```bash
make lint && make test && make build  # All must pass
# Verify test coverage ≥ 80%
go tool cover -func=coverage.out | tail -1  # Must show ≥80.0%
git add . && git commit -m "feat: ..."  # Conventional Commits format
```

**⚠️ Test Coverage Requirement (ENFORCED):**
- **MUST achieve ≥80% test coverage** before committing code
- Use `go tool cover -func=coverage.out` to check coverage
- If coverage is below 80%, **MUST** add more test cases
- Focus on 0% coverage functions and methods first

### 5. PR Management

#### Creating a PR
```bash
# When ready to create PR
git push -u origin feat/<issue-id>-<short-name>
gh pr create --title "feat: description" --body "Fixes #<issue-id>

## Changes
- Change 1
- Change 2

## Testing
- Test coverage: X%
- Manual testing completed"
```

#### Reviewing a PR
```bash
# When user says: "审查 PR #123" or "Review PR #123"
gh pr view 123                    # View PR details
gh pr checks 123                  # Check CI status
gh pr diff 123                    # Review code changes
gh pr review 123 --approve        # Approve if good
gh pr merge 123 --squash          # Merge if approved
```

**Consequence of non-compliance: PR will be rejected and must restart process.**

## Development Commands

- **Development**: `make dev` (starts Gin server with hot reload)
- **Build**: `make build` (compiles Go binary)
- **Test**: `make test` (runs all tests with coverage)
- **Lint**: `make lint` (runs golangci-lint + go fmt + go vet)
- **Database**: `make db-migrate` (run database migrations)
- **Health Check**: GET `/api/health`

## Environment Setup

Required environment variables in `.env`:
```
# Server Configuration
GIN_MODE=debug
PORT=8080
HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=vending_machine
DB_PASSWORD=your_password
DB_NAME=vending_machine_dev
DB_NAME_TEST=vending_machine_test

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_change_this_in_production
JWT_EXPIRES_HOURS=24

# WeChat Configuration
WECHAT_APP_ID=your_wechat_app_id
WECHAT_APP_SECRET=your_wechat_app_secret

# WeChat Pay Configuration
WECHAT_PAY_MERCHANT_ID=your_merchant_id
WECHAT_PAY_API_KEY=your_api_key
WECHAT_PAY_NOTIFY_URL=https://yourdomain.com/api/callback/wechat

# MQTT Device Communication
MQTT_BROKER=tcp://localhost:1883
MQTT_USERNAME=mqtt_user
MQTT_PASSWORD=mqtt_password
MQTT_CLIENT_ID=vending_machine_server
```

Optional for development:
```
MOCK_MODE=false  # Enable mock data responses
LOG_LEVEL=debug # Set logging level
REDIS_URL=redis://localhost:6379
OPENTELEMETRY_ENDPOINT=http://localhost:4317
```

## Core Architecture

This is a smart vending machine platform using Gin framework with the following structure:

### API Layer (`internal/handlers/`)
**Account Management:**
- `POST /api/Account/WeChatLogin` - WeChat user login
- `GET /api/Account/CheckLogin` - Check login status
- `GET /api/Account/GetUserInfo` - Get user information
- `GET /api/Account/CheckUserInfo` - Check user info by code

**Member Management:**
- `POST /api/Member/Update` - Update member information
- `POST /api/Member/AddFranchiseIntention` - Add franchise intention

**Machine Management:**
- `POST /api/Machine/GetPaging` - Get paginated machine list
- `GET /api/Machine/GetList` - Get machine list
- `GET /api/Machine/Get` - Get machine details
- `GET /api/Machine/GetProductList` - Get machine products
- `GET /api/Machine/OpenOrClose` - Toggle business status

**Order Management:**
- `POST /api/Order/GetPaging` - Get paginated orders
- `GET /api/Order/Get` - Get order details
- `POST /api/Order/Create` - Create new order
- `POST /api/Order/Refund` - Process refund

**Payment:**
- `GET /api/Payment/Get` - Get payment info
- `GET /api/Payment/Query` - Query payment status
- `POST /api/Callback/PaymentResult` - WeChat payment callback

All APIs use MySQL for data persistence with WeChat integration.

### Contract-First Development (`internal/contracts/`)
All API requests/responses are validated using Go structs:
- `AccountRequest` → `AccountResponse` (Login/Auth)
- `MemberRequest` → `MemberResponse` (Member management)
- `MachineRequest` → `MachineResponse` (Machine operations)
- `OrderRequest` → `OrderResponse` (Order processing)
- `PaymentRequest` → `PaymentResponse` (Payment processing)
- Breaking changes to contracts require PR and README changelog update

### Key Processing Pipeline
1. **Authentication**: JWT middleware for protected routes
2. **WeChat Integration**: WeChat login and payment processing
3. **Request Validation**: Gin middleware validates incoming requests
4. **Business Logic**: Service layer processes vending machine operations
5. **MQTT Communication**: Real-time device communication
6. **Data Persistence**: Repository layer handles database operations
7. **Response Formation**: Structured JSON responses with error handling

## Testing and Development

**IMPORTANT: All development work MUST follow the above standard workflow, from Issue to PR completion.**

- **Mock Mode**: Set `MOCK_MODE=true` for testing without database
- **Test Coverage**: Maintain ≥80% test coverage across all packages (ENFORCED)
- **Task Entry Point**: ONLY start from GitHub Issues (labels: `backend`/`api`/`docs`)
- **Commit Standards**: Strictly follow Conventional Commits format
- **Issue Linking**: PRs MUST include `Fixes #<issue-id>`

### Development Checklist

- [ ] **切换主分支并拉取最新代码** (`git checkout main && git pull`)
- [ ] **验证工作目录干净** (`git status`)
- [ ] Review and understand Issue requirements
- [ ] Check dependencies using `docs/WORKFLOWS/TASK_PARALLEL.md` if complex task
- [ ] Create feature branch (基于最新main分支)
- [ ] Use TodoWrite to plan development tasks
- [ ] Implement with real-time progress updates
- [ ] Final quality checks (lint + test + build)
- [ ] Commit code and create PR with `Fixes #<issue-id>`

## Code Patterns

### Type Safety
- Import types from `internal/contracts`
- Use strong typing with Go structs
- Path alias mapping: `internal/*` for internal packages

### API Response Handling
All responses include optional `Meta.Warnings[]` for user feedback:
```go
type APIResponse struct {
    Data interface{} `json:"data"`
    Meta *Meta      `json:"meta,omitempty"`
}

type Meta struct {
    Warnings []string `json:"warnings,omitempty"`
    Count    int      `json:"count,omitempty"`
}
```

### Error Handling
Use structured error handling:
```go
import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, err error, code int) {
    c.JSON(code, gin.H{
        "error": err.Error(),
        "code":  code,
    })
}
```

## Important Constraints

- **Database Transactions**: Use transactions for multi-table operations
- **Input Validation**: Validate all inputs using struct tags and custom validators
- **Authentication**: JWT-based authentication for protected endpoints
- **Rate Limiting**: 100 requests per minute per IP on write operations

## Component Structure

- `cmd/server/` - Main application entry point
- `internal/handlers/` - HTTP request handlers (Account, Member, Machine, Order, Payment)
- `internal/services/` - Business logic layer (JWT, Cache, Device communication)
- `internal/repositories/` - Database access layer
- `internal/models/` - Data models and entities (Member, Machine, Order, Product)
- `internal/contracts/` - API contract definitions
- `internal/middleware/` - Gin middleware components (Auth, CORS, Logger)
- `internal/config/` - Configuration management (Database, WeChat)
- `pkg/wechat/` - WeChat SDK integration

## Mock Development

For API development without database dependencies:
- Set `MOCK_MODE=true` in environment
- Mock data structure matches production schemas
- Supports all CRUD operations with in-memory storage

## GitHub Operations

### Essential Commands

#### Task Management
```bash
# View issue details
gh issue view <issue-id>
gh issue edit <issue-id> --add-label "in-progress"

# Create Epic or requirement
gh issue create --title "[Epic] Feature Name" --label "epic,backend"

# Track progress
gh issue list --label "epic" --state open
gh issue list --label "blocked" --state open
```

#### PR Management
```bash
# Create PR
gh pr create --title "feat: feature name" --body "Fixes #<issue-id>"

# Review PR
gh pr view 123
gh pr checks 123
gh pr diff 123
gh pr review 123 --approve
gh pr merge 123 --squash
```

### PR Merge Standards
Follow validation flow from `docs/WORKFLOWS/PR_MERGE.md`:
- **Required checks**: CI/CD status, merge conflicts, Issue linking, functional validation
- **Risk classification**: Low risk (auto-merge), Medium risk (extra validation), High risk (human review)
- **Emergency fixes**: Tagged `urgent`/`hotfix` can bypass certain checks

### Documentation Structure
```
docs/
├── QUICK_START.md              # 🚀 1-minute getting started
├── README.md                   # 📋 Documentation index
├── WORKFLOWS/                  # 📚 Core workflows
│   ├── DEVELOPMENT.md         # Development workflow
│   ├── REQUIREMENT.md         # Requirement analysis
│   ├── PR_MERGE.md           # PR review & merge
│   └── TASK_PARALLEL.md      # Parallel task handling
├── EMERGENCY/                  # 🚨 Emergency procedures
│   ├── HOTFIX.md             # Hotfix process
│   └── ROLLBACK.md           # Rollback guide
├── PATTERNS/                   # 🔧 Common patterns
│   ├── API_DESIGN.md         # API design patterns
│   └── ERROR_HANDLING.md     # Error handling
├── PRD/                        # 📄 Product requirements
└── archive/                    # 📦 Archived guides
```

## Database Configuration

### Technology Stack
- **Database**: MySQL 8.0+
- **Migration Tool**: GORM Auto Migration
- **ORM**: GORM
- **Authentication**: JWT + WeChat Login
- **Payment**: WeChat Pay API
- **Device Communication**: MQTT Protocol
- **Core tables**: members, machines, machine_owners, products, machine_products, orders, franchise_intentions

### Development Commands
- **Database Migration**: `make db-migrate` (apply pending migrations)
- **Database Rollback**: `make db-rollback` (rollback last migration)
- **Database Reset**: `make db-reset` (drop and recreate database)
- **Database Seed**: `make db-seed` (populate with sample data)

### Environment Setup
```bash
# Development (.env)
DB_HOST=localhost
DB_PORT=3306
DB_USER=vending_machine
DB_PASSWORD=password
DB_NAME=vending_machine_dev

# Test (.env.test)
DB_NAME=vending_machine_test
```

## Deployment Configuration

### Environment Comparison
| Feature | Development | Production |
|---------|-------------|------------|
| Database | Local MySQL | Production MySQL cluster |
| Logging | Debug level | Info level |
| Auth | Relaxed JWT | Strict JWT + HTTPS |
| Purpose | Development & testing | Live service |
| Auto Deploy | Manual | CI/CD pipeline |

### Deployment Commands
- **Development**: `make dev` (starts with hot reload)
- **Production Build**: `make build-prod` (optimized binary)
- **Health Check**: `curl http://localhost:8080/api/health`

### Required Environment Variables
**Development** (set in .env):
```bash
GIN_MODE=debug
LOG_LEVEL=debug
MOCK_MODE=false
```

**Production** (set in deployment):
```bash
GIN_MODE=release
LOG_LEVEL=info
DB_HOST=production-mysql-host
JWT_SECRET=production-secret
```

### Deployment Verification
```bash
# Health check
GET /api/health

# Database connectivity test
GET /api/health/db

# Authentication test
POST /api/auth/login
```

## Quality Assurance

### Automated Validation Flow
Based on `docs/AGENT_PR_MERGE_GUIDE.md` standards:

**Required Checks** (must pass before merge):
1. **CI/CD Status**: All GitHub Actions checks pass
2. **Merge Conflicts**: No conflicts with target branch
3. **Issue Linking**: PR has valid `Fixes #<issue-id>` reference
4. **Functional Validation**: Core APIs respond correctly

### Validation Commands
```bash
# Code quality checks
make pre-commit    # Runs lint, test, build

# Functional testing
make health-check  # Validates /api/health endpoint
make test-api      # Tests all API endpoints

# Full validation pipeline
make deploy-check  # Complete pre-deployment validation
```

### Risk Classification
- **🟢 Low Risk** (auto-merge): `docs:`, `style:`, `test:`, `chore:` changes
- **🟡 Medium Risk** (extra validation): `feat:`, `fix:`, `refactor:` changes  
- **🔴 High Risk** (human review): `internal/contracts/*`, database schema, security changes

### Emergency Procedures
- **Urgent fixes**: Tag with `urgent`/`hotfix` labels for expedited merge
- **Rollback**: Use `git revert <commit-hash>` + immediate notification
- **Issue escalation**: Add `needs-human-review` label for complex problems

## 🚀 Quick Reference for Claude Agents

### Common Task Scenarios

#### Scenario 1: User wants to discuss a feature
```
User: "我想讨论一下订单退款功能"
Claude Actions:
1. Read docs/WORKFLOWS/REQUIREMENT.md
2. Analyze business requirements
3. Create PRD in docs/PRD/order-refund.md
4. Define acceptance criteria
5. Create issue: gh issue create --label "backend"
```

#### Scenario 2: User wants to start development
```
User: "处理 #54"
Claude Actions:
1. Read docs/WORKFLOWS/DEVELOPMENT.md
2. Check issue: gh issue view 54
3. Create branch: git checkout -b feat/54-name
4. Use TodoWrite to plan tasks
5. Implement and test
6. Create PR when complete
```

#### Scenario 3: User wants parallel tasks
```
User: "这些任务可以并行吗？"
Claude Actions:
1. Read docs/WORKFLOWS/TASK_PARALLEL.md
2. Check file conflicts between tasks
3. Identify dependencies
4. Create parallel task groups
```

#### Scenario 4: Emergency fix needed
```
User: "生产环境出bug了！"
Claude Actions:
1. Read docs/EMERGENCY/HOTFIX.md
2. Create hotfix branch from main
3. Apply minimal fix
4. Create urgent PR with "hotfix" label
```

### Workflow Entry Points
| Task Type | Start Here | Purpose |
|-----------|------------|----------|
| 🚀 Quick Start | `docs/QUICK_START.md` | 1-minute guide to common tasks |
| 💻 Development | `docs/WORKFLOWS/DEVELOPMENT.md` | Complete dev workflow |
| 📋 Requirements | `docs/WORKFLOWS/REQUIREMENT.md` | Requirement analysis & PRD |
| 🔄 PR Review | `docs/WORKFLOWS/PR_MERGE.md` | PR validation & merge |
| ⚡ Parallel Tasks | `docs/WORKFLOWS/TASK_PARALLEL.md` | Dependency & parallel execution |

### Common GitHub Commands
```bash
# Check project status
gh issue list --label "epic" --state open
gh issue list --label "blocked" --state open

# Dependency checking (use with TASK_PARALLEL.md)
gh issue view <issue-id> --json body -q .body | grep -E "- \[ \] #[0-9]+"

# Quality gates (Dev Agent)  
make lint && make test && make build
go tool cover -func=coverage.out | tail -1  # Verify ≥80% coverage

# Standard PR creation
gh pr create --title "feat: description" --body "Fixes #<issue-id>"
```

### Documentation Navigation
- 🚀 **Quick Start**: `docs/QUICK_START.md` - 1-minute guide
- 📋 **Index**: `docs/README.md` - Complete documentation index
- 🔧 **Workflows**: `docs/WORKFLOWS/` - All core workflows
- 🚨 **Emergency**: `docs/EMERGENCY/` - Hotfix and rollback procedures
- 🎨 **Patterns**: `docs/PATTERNS/` - Common design patterns

### Success Indicators
- ✅ **Use correct workflow guide** for the task at hand
- ✅ **Follow streamlined processes** in WORKFLOWS/ directory
- ✅ **GitHub operations** using recommended gh commands
- ✅ **Quality standards** meeting all validation checkpoints
- ✅ **TodoWrite usage** for task planning and tracking

## 📌 Hook System Overview

The project uses automated hooks to ensure code quality at key checkpoints:

### Active Hooks
- **Branch Creation**: Validates creating from latest main branch
- **Commit**: Checks code quality and test coverage (≥80%)
- **PR Creation**: Ensures all quality standards are met
- **Contract Changes**: Reminds about API compatibility

### Hook Configuration
- Configuration: `.claude/settings.json`
- Scripts: `.claude/scripts/`
- No UserPromptSubmit hooks - won't interrupt normal conversation

---

**Follow the streamlined workflows for efficient development!** 🚀

*CLAUDE.md last updated: 2025-08-13*