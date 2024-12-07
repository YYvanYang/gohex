# Contributing to GoHex

Thank you for your interest in contributing to GoHex! This document will guide you through the contribution process.

## 开发流程

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的修改 (`git commit -m 'feat: add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建一个 Pull Request

## 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范，提交格式如下：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- feat: 新功能
- fix: 修复问题
- docs: 文档修改
- style: 代码格式修改
- refactor: 代码重构
- perf: 性能优化
- test: 测试相关
- chore: 其他修改

### Scope

可选的修改范围：

- auth: 认证相关
- user: 用户相关
- event: 事件相关
- cache: 缓存相关
- db: 数据库相关
- api: API 相关
- config: 配置相关
- deps: 依赖相关

### 示例

```
feat(user): add user registration

- Add user registration endpoint
- Add email verification
- Add password validation

Closes #123
```

## 代码规范

1. 遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
2. 所有代码必须通过 `golangci-lint` 检查
3. 所有公开的函数和类型必须有文档注释
4. 测试覆盖率必须达到 80% 以上

## 目录结构

```
.
├── cmd/                  # 应用程序入口
├── internal/            # 内部代码
│   ├── application/    # 应用层
│   ├── domain/        # 领域层
│   └── infrastructure/ # 基础设施层
├── pkg/                # 公共包
├── scripts/           # 脚本
└── test/              # 测试
```

## 开发环境设置

1. 安装依赖：
bash
go mod download
```

2. 安装开发工具：
```bash
make install-tools
```

3. 运行测试：
```bash
make test
```

4. 运行 lint：
```bash
make lint
```

## 测试规范

1. 单元测试：
   - 文件名格式：`xxx_test.go`
   - 使用 testify 包进行断言
   - 每个包必须有测试
   - 使用表驱动测试方式

2. 集成测试：
   - 位于 `test/integration` 目录
   - 使用 Docker Compose 进行依赖服务管理
   - 测试前需要确保依赖服务正常运行

## 文档规范

1. API 文档：
   - 使用 Swagger/OpenAPI 规范
   - 位于 `api/swagger` 目录
   - 每个 API 必须有详细的描述和示例

2. 架构文档：
   - 位于 `docs/architecture` 目录
   - 使用 PlantUML 绘制架构图
   - 包含系统设计决策和理由

## 发布流程

1. 版本号规范：遵循 [Semantic Versioning](https://semver.org/)
2. 每个版本必须有 CHANGELOG
3. 发布前必须通过所有测试
4. 使用 tag 标记版本：`git tag v1.0.0`

## 问题反馈

1. 使用 GitHub Issues 进行问题反馈
2. 使用提供的 issue 模板
3. 提供详细的复现步骤和环境信息

## 安全问题

如果您发现了安全漏洞，请不要直接在 Issues 中提出，而是发送邮件到 security@example.com。

## 许可证

通过提交 pull request，您同意将您的代码按照项目的开源协议进行授权。

## 行为准则

请阅读并遵守我们的 [行为准则](CODE_OF_CONDUCT.md)。

```
</```rewritten_file>