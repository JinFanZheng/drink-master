# GitHub Issue 自动化命令

本项目提供了一些自动化命令来帮助管理GitHub Issues。

## 🤖 可用命令

在任何Issue的评论中输入以下命令：

### `/claim`
- **功能**: 认领这个Issue并标记为进行中
- **效果**: 
  - 将你设为Issue的负责人
  - 添加`in-progress`标签
  - 发送确认消息

**示例**:
```
/claim
```

### `/unclaim` 
- **功能**: 取消认领这个Issue
- **效果**:
  - 移除你的负责人身份
  - 移除`in-progress`标签
  - 发送确认消息

**示例**:
```
/unclaim
```

### `/help`
- **功能**: 显示帮助信息和开发指南链接
- **效果**: 显示可用命令和项目文档链接

**示例**:
```
/help
```

## 🏷️ 自动标签

当创建新Issue时，系统会根据标题和内容自动添加标签：

### 优先级标签
- `priority-high`: 包含"urgent"、"critical"关键词
- `priority-medium`: 包含"important"关键词  
- `priority-low`: 默认优先级

### 类型标签
- `bug`: 包含"bug"、"error"、"fix"关键词
- `enhancement`: 包含"feature"、"enhancement"关键词
- `documentation`: 包含"doc"、"documentation"关键词
- `api`: 包含"api"关键词
- `backend`: 包含"backend"关键词

## 👋 新贡献者欢迎

首次提交Issue的贡献者会收到自动欢迎消息，包含：
- 项目欢迎信息
- Issue提交指南
- 开发文档链接

## 📋 开发工作流

1. **查看Issue**: 浏览可用的Issue
2. **认领Issue**: 使用`/claim`命令认领感兴趣的Issue
3. **开始开发**: 按照[开发指南](../docs/AGENT_ONBOARDING.md)进行开发
4. **提交PR**: 完成开发后提交Pull Request
5. **代码审查**: 等待代码审查和合并

## 🎯 注意事项

- 确保测试覆盖率达到80%以上
- 遵循[协作指南](../docs/ROLES_COLLABORATION.md)中的规范
- 使用conventional commit格式提交代码
- 一次只认领一个Issue，完成后再认领下一个