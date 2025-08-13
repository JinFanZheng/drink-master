# 需求分析工作流程

## 1. 需求分析

### 当用户说"我想要XXX功能"时

#### 理解需求
1. 分析核心诉求
2. 评估技术可行性
3. 确定功能边界

#### 创建PRD文档
```markdown
# docs/PRD/<feature-name>.md

## 需求背景
- 用户痛点
- 业务价值

## 功能描述
- 核心功能点
- 用户流程

## 验收标准
- [ ] 标准1
- [ ] 标准2
- [ ] 标准3

## 成功指标
- 指标1: 具体数值
- 指标2: 具体数值
```

#### 创建需求Issue
```bash
gh issue create \
  --title "[需求] <feature-name>" \
  --label "product,priority-high" \
  --body "## 需求描述
  
## 验收标准
- [ ] 标准1
- [ ] 标准2

## PRD文档
docs/PRD/<feature-name>.md"
```

## 2. 任务拆解

### 将需求转化为开发任务

#### 评估工作量
- 小任务（<1天）：直接创建开发Issue
- 中任务（1-3天）：拆分2-3个子任务
- 大任务（>3天）：创建Epic + 多个子任务

#### 创建开发任务
```bash
# 单个任务
gh issue create \
  --title "[开发] 实现<feature>" \
  --label "backend" \
  --body "## 任务描述

## 子任务
- [ ] 实现数据模型
- [ ] 实现API接口
- [ ] 添加测试用例

## 验收标准
- [ ] 测试覆盖率≥80%
- [ ] API文档更新

Related to #<product-issue-id>"
```

#### 创建Epic（复杂功能）
```bash
# Epic
gh issue create \
  --title "[Epic] <feature-name>" \
  --label "epic" \
  --body "## Epic描述

## 子任务列表
- [ ] #task1 - 任务1描述
- [ ] #task2 - 任务2描述
- [ ] #task3 - 任务3描述

## 依赖关系
- Task2 依赖 Task1
- Task3 可与Task2并行"
```

## 3. 任务优先级

### 优先级标签使用
- `priority-high`: 本周必须完成
- `priority-medium`: 本迭代内完成
- `priority-low`: 有时间再做

### 里程碑设置
```bash
# 将任务加入里程碑
gh issue edit <issue-id> --milestone "Sprint-2024-W45"
```

## 4. 需求变更处理

### 评估变更影响
1. 已开发功能的影响
2. 进行中任务的影响
3. 测试用例的更新

### 更新文档
1. 更新PRD文档
2. 在Issue中添加变更说明
3. 通知相关开发任务

## 常见需求模式

### CRUD功能
```
需求：管理XXX数据
拆解：
1. 数据模型设计
2. 创建API
3. 查询API（列表+详情）
4. 更新API
5. 删除API
```

### 集成功能
```
需求：集成XXX服务
拆解：
1. 研究API文档
2. 设计集成方案
3. 实现接口调用
4. 错误处理
5. 集成测试
```

### 优化功能
```
需求：优化XXX性能
拆解：
1. 性能分析
2. 确定优化点
3. 实施优化
4. 性能测试
5. 监控指标
```