# PM/Planner Agent Onboarding

本指南面向"项目管理与执行协调"的代理（PM agents）。

## 角色定位
**PM负责"执行层面"**：将已定义的产品需求转化为可执行的开发任务，协调资源和进度。

**与Product Agent的分工**：
- **Product**: 需求发现 → 方案设计 → PRD输出
- **PM**: 任务拆解 → 资源协调 → 进度跟踪 → 交付管理

## 目标
- 将产品需求高效转化为开发任务，确保按时保质交付

## 关键入口
- 仓库: https://github.com/ddteam/drink-master
- Roadmap 看板: https://github.com/users/ddteam/projects/1
- **开发流程**: 必须先阅读 `docs/AGENT_ONBOARDING.md` 了解强制开发流程
- 文档入口: `README.md`、`README.md（项目管理部分）`、`CLAUDE.md`
- 命令行工具：建议用 `gh` 批量管理 Issues/Milestones/Project（示例：`gh issue create|edit`、`gh api`、`gh project item-add`）

## ⚠️ 重要提醒
**所有创建的 Issue 都必须遵循 `docs/AGENT_ONBOARDING.md` 中定义的标准开发流程。**

## 任务模型
- Issue 是唯一任务实体；用标签表达维度：
  - 领域: `backend` / `api` / `docs` / `test`
  - 优先级: `priority-high|medium|low`
  - 状态: `in-progress` / `blocked`
  - 类型（可选）: `epic`（史诗/主题）
- 里程碑: M1/M2/M3（见 GitHub Milestones）
- 项目看板: 跟踪状态（Todo/In progress/Blocked/Review/Done）

## 迭代（Sprint）管理
- 周期：建议 1–2 周/迭代。每个迭代创建一个 Milestone（如 `Sprint-YYYY-WW`），所有本迭代 Issues 必须挂此里程碑。
- 计划：在迭代开始前 30–45 分钟确定目标（1–3 个可交付价值），拆分为 Epic/Issue 并指派里程碑与优先级。
- 执行：Agents 认领 `/claim` → `in-progress`，遇阻塞打 `blocked` 并评论；PM 每日关注 Blocked/Review。
- 评审：迭代末演示功能（Mock/真实环境均可），Ready for review 的 PR 自动将 Issue 标 `review`；按 DoD 验收。
- 完结：合并后自动 Done；里程碑关闭时生成 Release Note（建议将变更摘要追加到 README 或 docs/）。

## PM 核心职责

### 1. 任务拆解与规划
**输入**：Product Agent提供的PRD和功能需求
**输出**：具体的开发任务(Issues)和依赖关系

- **Epic管理**：基于产品需求创建Epic，管理子任务依赖
- **任务拆解**：将功能需求分解为`backend`/`api`/`docs`具体任务
- **依赖规划**：使用`docs/TASK_DEPENDENCY_PLANNING.md`构建任务依赖图
- **批次组织**：识别可并行执行的任务，优化开发效率

### 2. 资源协调与分派
- **容量评估**：评估开发团队容量，合理分配任务
- **技能匹配**：根据任务类型匹配合适的开发代理
- **冲突预防**：协调潜在的文件修改冲突（参考依赖规划指南）
- **阻塞处理**：快速响应`blocked`任务，协调解除阻塞

### 3. 进度跟踪与节奏控制
- **看板管理**：每日检查Blocked/Review列，确保流程顺畅
- **里程碑跟踪**：监控Sprint目标进度，及时预警风险
- **质量把关**：确保PR关联Issue，CI通过，符合DoD标准
- **发布协调**：管理发布计划，协调多个功能的集成测试

### 4. 交付与复盘
- **验收管理**：组织功能演示，确保符合产品预期
- **变更控制**：管理需求变更，评估对进度和资源的影响
- **过程改进**：收集开发过程中的问题，优化流程效率

## 状态自动化（已配置）
- `/claim` 评论: 自动指派 assignee 并添加 `in-progress` 标签
- 打标签: `in-progress`/`blocked` 会同步触发项目 Status 变更（In progress/Blocked）

## 例行检查清单
- [ ] 看板"Blocked"列中的任务是否已收到响应与行动
- [ ] 本周必须完成的里程碑任务是否按时推进
- [ ] 关键契约/接口变更是否同步更新文档与变更记录
- [ ] CI 是否稳定（main 分支绿灯）

## PM专用工具与模板

### Epic模板 (PM负责创建)
```markdown
# [Epic] 功能主题

## 📋 基本信息
**来源PRD**: docs/PRD/<topic>.md
**Product负责人**: @product-agent
**预估工期**: X周
**里程碑**: M1/M2/M3

## 🏗️ 任务依赖图
[参考 TASK_DEPENDENCY_PLANNING.md 构建]

## 📋 开发任务清单
### Batch 1: 基础设施
- [ ] #XX 基础任务A
- [ ] #XX 基础任务B

### Batch 2: 功能开发 (可并行)
- [ ] #XX 后端任务A ⭐
- [ ] #XX API任务B ⭐
- [ ] #XX 文档任务C ⭐

## 🎯 交付计划
**Sprint目标**: 明确的可演示价值
**风险预案**: 关键风险和缓解措施
**资源需求**: 所需的开发代理类型和数量
```

### PM日常操作清单
- [ ] 每日看板巡检：处理Blocked和Review任务
- [ ] 每周里程碑检查：确保Sprint目标可达成
- [ ] Epic进度跟踪：更新任务完成状态
- [ ] 依赖冲突协调：提前识别和解决文件冲突
- [ ] 发布计划管理：协调多功能的集成测试

## 参考文档
- **上游依赖**: `docs/PRODUCT_ONBOARDING.md` - 了解Product输入
- **依赖规划**: `docs/TASK_DEPENDENCY_PLANNING.md` - DAG任务规划方法
- **开发流程**: `docs/AGENT_ONBOARDING.md` - 开发代理标准流程
- **技术规范**: `CONTRIBUTING.md` - 代码提交和PR规范