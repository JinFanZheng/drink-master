# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ⚠️ Agent Collaboration Framework (REQUIRED)

**All Claude agents MUST understand their role and follow the collaboration framework:**

### 🎯 Role Identification
Before starting ANY task, identify your role based on the work being performed:
- **Product Agent**: User research, PRD creation, feature validation (🎯)
- **PM Agent**: Epic management, task coordination, dependency planning (📊)  
- **Dev Agent**: Code development, technical implementation (💻)

### 📚 Required Reading (Role-Based)
**ALL Agents** must first read:
1. `docs/ROLES_COLLABORATION.md` - Overall collaboration framework
2. Your role-specific onboarding document:
   - Product → `docs/PRODUCT_ONBOARDING.md`
   - PM → `docs/PM_ONBOARDING.md`
   - Dev → `docs/AGENT_ONBOARDING.md`

### 1. Pre-Development Setup (MANDATORY)
```bash
# Execute in order (Dev Agent workflow):
git checkout main && git pull origin main           # Ensure latest code
git status                                          # Confirm clean working directory
gh issue view <issue-id>                           # Review task requirements
gh issue edit <issue-id> --add-label "in-progress" # Mark as in-progress
git checkout -b feat/<issue-id>-<short-name>       # Create feature branch
make lint && make test && make build               # Basic quality checks
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

### 5. PR Creation (MANDATORY)
```bash
git push -u origin feat/<issue-id>-<short-name>
gh pr create --title "..." --body "Fixes #<issue-id> ..."
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

### Role-Specific Workflow Checklists

#### 🎯 Product Agent Checklist
- [ ] Read `docs/PRODUCT_ONBOARDING.md` for detailed workflow
- [ ] Create product requirement Issues with `product` label
- [ ] Output PRD documents to `docs/PRD/<topic>.md`
- [ ] Define clear DoD and success metrics
- [ ] Participate in feature validation and acceptance

#### 📊 PM Agent Checklist  
- [ ] Read `docs/PM_ONBOARDING.md` for detailed workflow
- [ ] Create Epic Issues based on Product PRDs
- [ ] Use `docs/TASK_DEPENDENCY_PLANNING.md` for dependency analysis
- [ ] Break down Epics into specific development tasks
- [ ] Coordinate dev resources and track progress

#### 💻 Dev Agent Checklist
- [ ] Read `docs/AGENT_ONBOARDING.md` for detailed workflow
- [ ] **切换主分支并拉取最新代码** (`git checkout main && git pull`)
- [ ] **验证工作目录干净** (`git status`)
- [ ] Review and understand Issue requirements
- [ ] Check dependencies using `docs/TASK_DEPENDENCY_PLANNING.md`
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

## Team Collaboration Guidelines

### Agent Collaboration Framework
**Follow the structured collaboration model defined in `docs/ROLES_COLLABORATION.md`:**

#### 🔄 Standard Workflow
1. **Product Agent** → User research → PRD creation → Success metrics
2. **PM Agent** → Epic creation → Task breakdown → Dependency planning  
3. **Dev Agent** → Code implementation → Quality assurance → PR creation
4. **All Agents** → Feature validation → Release coordination → Data analysis

#### 🚫 Collaboration Boundaries
- **Don't duplicate**: Each role has specific responsibilities, avoid overlap
- **Don't skip steps**: Follow the sequential workflow stages
- **Don't work in isolation**: Use designated communication and handoff points

### GitHub Operations (gh commands)
**Essential commands for each role:**

#### Product Agent Commands
```bash
# Create product requirement
gh issue create --title "[Product] Feature Name" --label "product,priority-high" --body-file prd-template.md

# Validate completion  
gh issue comment <issue-id> --body "✅ Product validation passed"
```

#### PM Agent Commands  
```bash
# Create Epic from PRD
gh issue create --title "[Epic] Feature Name" --label "epic,backend" --milestone "M1"

# Track progress
gh issue list --label "epic" --state open
gh project item-list <project-id>
```

#### Dev Agent Commands
```bash
# Start development (following AGENT_ONBOARDING.md)
gh issue view <issue-id>
gh issue edit <issue-id> --add-label "in-progress"

# Create PR
gh pr create --title "feat: feature name" --body "Fixes #<issue-id>"
```

### PR Merge Standards
Follow automated validation flow from `docs/AGENT_PR_MERGE_GUIDE.md`:
- **Required checks**: CI/CD status, merge conflicts, Issue linking, functional validation
- **Risk classification**: Low risk (auto-merge), Medium risk (extra validation), High risk (human review)
- **Emergency fixes**: Tagged `urgent`/`hotfix` can bypass certain checks

### Documentation Structure
All agents should understand the documentation organization:
```
docs/
├── README.md                    # 📋 Documentation index
├── ROLES_COLLABORATION.md       # 🎯 Core collaboration guide  
├── [ROLE]_ONBOARDING.md        # 📚 Role-specific workflows
├── PRD/                        # 📄 Product requirements
├── Sprint/                     # 📊 Sprint planning
├── Guides/                     # 📚 Technical guides
└── Operations/                 # ⚙️ Ops documentation
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

### Role-Based Entry Points
| Role | Start Here | Key Documents | Main Output |
|------|------------|---------------|-------------|
| 🎯 Product | `docs/PRODUCT_ONBOARDING.md` | User research, PRD templates | `docs/PRD/*.md` |
| 📊 PM | `docs/PM_ONBOARDING.md` | Epic management, dependency planning | GitHub Epics + Issues |
| 💻 Dev | `docs/AGENT_ONBOARDING.md` | Code quality, branch workflow | Code + PRs |

### Common GitHub Commands
```bash
# Check project status
gh issue list --label "epic" --state open
gh issue list --label "blocked" --state open

# Dependency checking (use with TASK_DEPENDENCY_PLANNING.md)
gh issue view <issue-id> --json body -q .body | grep -E "- \[ \] #[0-9]+"

# Quality gates (Dev Agent)  
make lint && make test && make build
go tool cover -func=coverage.out | tail -1  # Verify ≥80% coverage

# Standard PR creation
gh pr create --title "feat: description" --body "Fixes #<issue-id>"
```

### Documentation Navigation
- 🎯 **Start**: `docs/README.md` - Complete documentation index
- 🤝 **Collaboration**: `docs/ROLES_COLLABORATION.md` - Role boundaries and workflows  
- 🔧 **Methods**: `docs/TASK_DEPENDENCY_PLANNING.md` - DAG dependency analysis
- 📋 **Process**: `docs/AGENT_PR_MERGE_GUIDE.md` - PR review and merge standards

### Success Indicators
- ✅ **Clear role identification** before starting any task
- ✅ **Proper document structure** following the new organization
- ✅ **GitHub operations** using recommended gh commands
- ✅ **Quality standards** meeting all validation checkpoints
- ✅ **Collaboration boundaries** respecting role-specific responsibilities

---

**Follow the Agent Collaboration Framework for efficient teamwork!** 🤝

*CLAUDE.md last updated: 2025-08-12*