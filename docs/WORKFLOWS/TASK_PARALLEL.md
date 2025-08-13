# 并行任务管理

## 快速判断：任务是否可并行？

### ✅ 可以并行的情况
- 不同的handler文件（如：user_handler.go vs product_handler.go）
- 不同的service文件
- 独立的功能模块
- 不同的API端点

### ❌ 必须串行的情况
- 修改相同的contract文件
- 修改相同的model文件
- B功能依赖A功能的输出
- 数据库schema有依赖关系

## 文件冲突检查

### 检查两个任务是否会冲突
```bash
# 假设有两个分支
git diff main...feat/task-1 --name-only > files1.txt
git diff main...feat/task-2 --name-only > files2.txt

# 查看是否有相同文件
comm -12 <(sort files1.txt) <(sort files2.txt)
# 如果有输出，说明有文件冲突，需要串行开发
```

## 创建并行任务组

### 示例：用户管理模块
```bash
# 1. 创建Epic
gh issue create \
  --title "[Epic] 用户管理模块" \
  --label "epic" \
  --body "## 功能列表
- 用户登录
- 用户注册  
- 密码重置
- 个人信息修改"

# 2. 创建可并行的任务（登录、注册、密码重置可并行）
gh issue create --title "实现用户登录API" --label "backend" \
  --body "Part of #<epic-id>
可与注册、密码重置并行开发"

gh issue create --title "实现用户注册API" --label "backend" \
  --body "Part of #<epic-id>
可与登录、密码重置并行开发"

gh issue create --title "实现密码重置API" --label "backend" \
  --body "Part of #<epic-id>
可与登录、注册并行开发"

# 3. 创建有依赖的任务
gh issue create --title "实现个人信息修改API" --label "backend" \
  --body "Part of #<epic-id>
**依赖**: #<login-issue-id> (需要用户登录后才能修改信息)"
```

## 并行开发实践

### 1. 分配任务给不同分支
```bash
# 开发者A
git checkout -b feat/101-user-login

# 开发者B（可同时进行）
git checkout -b feat/102-user-register

# 开发者C（可同时进行）
git checkout -b feat/103-password-reset
```

### 2. 定期同步main分支
```bash
# 每天开始工作前
git fetch origin main
git rebase origin/main
```

### 3. 及时合并完成的任务
先完成的任务先合并，减少后续冲突

## 常见并行模式

### 模式1：CRUD并行
```
商品管理：
├── Task A: 商品Model + Migration（先做）
├── Task B: 创建商品API（依赖A）
├── Task C: 查询商品API（依赖A）  
├── Task D: 更新商品API（依赖A）
└── Task E: 删除商品API（依赖A）

B/C/D/E 可以并行开发
```

### 模式2：模块并行
```
电商系统：
├── Module A: 用户模块（独立）
├── Module B: 商品模块（独立）
├── Module C: 购物车模块（依赖A和B）
└── Module D: 订单模块（依赖C）

A和B可以并行，C等待A和B，D等待C
```

### 模式3：前后端并行
```
功能开发：
├── Task A: API契约定义（先做）
├── Task B: 后端API实现（依赖A）
└── Task C: 前端界面开发（依赖A）

B和C可以并行（基于契约开发）
```

## 依赖标记规范

在Issue描述中标记依赖：
```markdown
## 依赖关系
- **前置依赖**: #123 (必须先完成)
- **软依赖**: #124 (最好先完成，但可并行)
- **被依赖**: #125, #126 (这些任务等待本任务)
```

## 并行冲突解决

### 预防措施
1. 开发前检查文件冲突
2. 合理划分模块边界
3. 及时沟通和同步

### 冲突处理
```bash
# 如果出现冲突
git rebase origin/main
# 解决冲突
git add .
git rebase --continue

# 或者使用merge（保留历史）
git merge origin/main
# 解决冲突
git add .
git commit
```