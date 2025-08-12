# Git Hooks - Pre-commit 检查

这个目录包含了项目的Git hooks，用于在提交前进行代码质量检查。

## 设置

项目已经自动配置Git hooks路径，如果你需要手动设置：

```bash
git config core.hooksPath .githooks
```

## Pre-commit Hook

`pre-commit` hook会在每次提交前自动运行以下检查：

### 检查项目
1. **go mod tidy** - 确保模块依赖整洁
2. **go fmt** - 检查代码格式
3. **go vet** - 静态代码分析
4. **golangci-lint** - 高级代码检查（如果已安装）
5. **go test** - 运行单元测试
6. **测试覆盖率** - 确保覆盖率≥78%

### 安装 golangci-lint

为了获得最佳的代码检查体验，建议安装 golangci-lint：

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 手动运行检查

你也可以使用 Makefile 手动运行相同的检查：

```bash
# 运行完整的pre-commit检查
make pre-commit

# 或者分别运行
make lint    # 代码检查
make test    # 运行测试
make build   # 编译检查
```

### 跳过检查

如果确实需要跳过pre-commit检查（不推荐），可以使用：

```bash
git commit --no-verify -m "你的提交信息"
```

### 故障排除

如果pre-commit检查失败：

1. **格式问题**: 运行 `go fmt ./...`
2. **测试失败**: 运行 `go test ./...` 查看详细错误
3. **覆盖率不足**: 运行 `go test -coverprofile=coverage.out ./...` 检查覆盖率
4. **模块依赖**: 运行 `go mod tidy`

## CI/CD 一致性

这些检查与GitHub Actions CI管道保持一致，确保本地开发和CI环境的一致性。