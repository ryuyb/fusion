# Fusion 直播聚合平台 - 项目结构

## 目录组织

### 根目录结构

```
fusion/
├── .spec-workflow/               # Spec 工作流程文档
│   ├── steering/                # Steering 文档
│   │   ├── product.md          # 产品概述
│   │   ├── tech.md             # 技术架构
│   │   └── structure.md        # 项目结构（本文件）
│   ├── templates/              # Spec 模板
│   ├── specs/                  # 规范文档
│   └── user-templates/         # 自定义模板
├── .github/                     # GitHub 配置
│   └── workflows/              # GitHub Actions 工作流
├── cmd/                         # 应用入口
│   └── server/                 # Server 入口
│       └── main.go             # 主程序入口点
├── configs/                     # 配置文件
│   ├── config.yaml             # 默认配置
│   ├── config.dev.yaml         # 开发环境配置
│   └── config.prod.yaml        # 生产环境配置
├── docs/                        # 文档目录
│   └── api/                    # API 文档
├── internal/                    # 内部代码（不对外暴露）
│   ├── app/                    # 应用层（依赖注入配置）
│   ├── application/            # 应用服务层
│   ├── domain/                 # 领域层（业务逻辑核心）
│   ├── infrastructure/         # 基础设施层
│   ├── interface/              # 接口层（HTTP handler）
│   └── pkg/                    # 共享包
├── scripts/                     # 脚本文件
│   ├── build-multiarch.sh      # 多架构构建脚本
│   └── migration/              # 迁移脚本
├── test/                        # 测试文件
│   └── integration/            # 集成测试
├── .air.toml                    # Air 热重载配置
├── .dockerignore                # Docker 忽略文件
├── .env.example                 # 环境变量示例
├── .gitignore                   # Git 忽略文件
├── Dockerfile                   # Docker 构建文件
├── docker-compose.yml           # Docker Compose 配置
├── go.mod                       # Go 模块文件
├── go.sum                       # Go 依赖校验
├── Makefile                     # Make 构建脚本
└── CLAUDE.md                    # 项目说明文档
```

### Clean Architecture 核心结构

```
internal/
├── app/                          # 依赖注入和模块组织
│   ├── module.go                # 应用模块入口，统一依赖注入
│   └── fx.go                    # FX 模块配置
│
├── domain/                       # 领域层（业务逻辑核心）
│   ├── entity/                  # 领域实体
│   │   └── user.go             # 用户实体
│   ├── repository/              # 仓储接口（抽象）
│   │   └── user.go             # 用户仓储接口
│   ├── service/                 # 领域服务
│   │   ├── user.go             # 用户领域服务
│   │   └── auth.go             # 认证领域服务
│   └── module.go                # 领域层 FX 模块
│
├── application/                  # 应用层（用例和应用服务）
│   ├── dto/                     # 数据传输对象
│   │   ├── request/            # 请求 DTO
│   │   │   └── user.go        # 用户请求
│   │   └── response/           # 响应 DTO
│   │       └── user.go        # 用户响应
│   ├── service/                 # 应用服务
│   │   ├── user_service.go     # 用户应用服务
│   │   ├── auth_service.go     # 认证应用服务
│   │   └── mocks/              # Mock（测试用）
│   └── module.go                # 应用层 FX 模块
│
├── infrastructure/               # 基础设施层
│   ├── database/                # 数据库相关
│   │   ├── ent/                # EntGO 生成代码
│   │   ├── schema/             # EntGO 模式定义
│   │   │   └── user.go        # 用户模式
│   │   └── module.go          # 数据库 FX 模块
│   ├── repository/              # 仓储实现
│   │   ├── user_repository.go  # 用户仓储实现
│   │   └── module.go          # 仓储 FX 模块
│   ├── config/                  # 配置管理
│   │   └── module.go          # 配置 FX 模块
│   ├── logger/                  # 日志系统
│   │   └── module.go          # 日志 FX 模块
│   └── module.go                # 基础设施 FX 模块
│
├── interface/                    # 接口层（HTTP）
│   ├── http/                    # HTTP 处理器
│   │   ├── dto/                # HTTP DTO
│   │   │   ├── request/       # HTTP 请求 DTO
│   │   │   └── response/      # HTTP 响应 DTO
│   │   ├── handler/            # HTTP 处理器
│   │   │   ├── user.go       # 用户处理器
│   │   │   └── auth.go       # 认证处理器
│   │   ├── middleware/         # 中间件
│   │   │   ├── auth.go       # 认证中间件
│   │   │   ├── cors.go       # CORS 中间件
│   │   │   ├── error_handler.go # 错误处理
│   │   │   ├── logger.go     # 日志中间件
│   │   │   └── recovery.go   # 恢复中间件
│   │   ├── route/              # 路由定义
│   │   │   ├── user.go       # 用户路由
│   │   │   ├── auth.go       # 认证路由
│   │   │   └── module.go     # 路由 FX 模块
│   │   └── module.go          # 接口层 FX 模块
│   └── module.go                # 接口层 FX 模块
│
└── pkg/                         # 共享包
    ├── auth/                   # 认证工具
    │   ├── jwt.go             # JWT 实现
    │   └── module.go          # 认证 FX 模块
    ├── errors/                 # 错误处理
    │   ├── errors.go          # 错误定义
    │   └── module.go          # 错误 FX 模块
    ├── utils/                  # 工具函数
    │   ├── password.go        # 密码工具
    │   └── module.go          # 工具 FX 模块
    ├── validator/              # 验证器
    │   └── validator.go       # 验证实现
    └── module.go               # 共享包 FX 模块
```

## 命名规范

### 文件命名
- **Go 源文件**: `snake_case.go`
  - 实体文件: `user.go`, `subscription.go`
  - 服务文件: `user_service.go`, `auth_service.go`
  - 仓储文件: `user_repository.go`
  - 处理器文件: `user_handler.go`
  - 路由文件: `user_route.go`
  - 模块文件: `module.go` (每个层级都有)

- **配置文件**: `snake_case.yaml`
  - `config.yaml`
  - `config.dev.yaml`
  - `config.prod.yaml`

- **测试文件**: `snake_case_test.go`
  - `user_service_test.go`
  - `user_repository_test.go`

### 代码命名

#### 变量命名
- **本地变量**: `camelCase`
  ```go
  userID := 123
  isActive := true
  ```

- **全局变量**: `PascalCase`（少数使用）
  ```go
  const DefaultPort = 8080
  ```

- **常量**: `UPPER_SNAKE_CASE`
  ```go
  const JWT_SECRET_KEY = "secret"
  const MAX_RETRY_COUNT = 3
  ```

- **结构体**: `PascalCase`
  ```go
  type UserService struct {
      repo     repository.UserRepository
      logger   *zap.Logger
  }

  type User struct {
      ID       int    `json:"id"`
      Username string `json:"username"`
      Email    string `json:"email"`
  }
  ```

- **接口**: `PascalCase`
  ```go
  type UserRepository interface {
      Create(ctx context.Context, user *User) error
      FindByID(ctx context.Context, id int) (*User, error)
  }
  ```

- **函数/方法**: `camelCase`
  ```go
  func (s *UserService) CreateUser(ctx context.Context, dto *CreateUserDTO) error {
      // ...
  }

  func (r *UserRepository) Create(ctx context.Context, user *User) error {
      // ...
  }
  ```

### 目录命名
- **统一使用单数形式**
  - `domain/` (不是 domains/)
  - `application/` (不是 applications/)
  - `infrastructure/` (不是 infrastructures/)
  - `interface/` (不是 interfaces/)

## 模块组织模式

### FX Module 结构
每个层级都有自己的 `module.go` 文件，提供 FX 依赖注入：

```go
// internal/domain/module.go
package domain

import "go.uber.org/fx"

var Module = fx.Module(
    "domain",
    // 提供领域层的依赖
)
```

### 依赖注入模式
- **构造函数注入**: 通过构造函数注入依赖
  ```go
  func NewUserService(
      repo repository.UserRepository,
      logger *zap.Logger,
  ) *UserService {
      return &UserService{
          repo:   repo,
          logger: logger,
      }
  }
  ```

- **接口依赖**: 依赖接口而非具体实现
  ```go
  type UserService struct {
      repo repository.UserRepository  // 接口
      logger *zap.Logger
  }
  ```

## 模块边界

### 层间依赖规则
```
Interface (HTTP)
    ↓
Application (Services)
    ↓
Domain (Entities & Interfaces)
    ↑
Infrastructure (Implementations)
```

### 关键边界
1. **Domain 独立**: Domain 层不依赖任何其他层
2. **Infrastructure 依赖 Domain**: 只依赖接口
3. **Application 依赖 Domain**: 只依赖接口和实体
4. **Interface 只依赖 Application**: 不直接调用 Infrastructure

### 模块通信
- **入站数据流**: Interface → Application → Domain
- **出站数据流**: Domain → Infrastructure
- **跨层通信**: 只能通过接口进行

## 代码组织原则

### 1. 单职责原则
每个文件都有明确的目的：
- `user_service.go`: 只包含用户相关的应用服务
- `user_repository.go`: 只包含用户相关的仓储逻辑
- `user_handler.go`: 只包含用户相关的 HTTP 处理器

### 2. 分层组织
按层组织代码，而不是按功能：
```
正确:
internal/
├── application/service/
│   ├── user_service.go
│   ├── auth_service.go
└── domain/entity/
    ├── user.go
    └── auth.go

错误（按功能组织）:
internal/
├── user/
│   ├── service.go
│   ├── repository.go
│   └── handler.go
└── auth/
    ├── service.go
    ├── repository.go
    └── handler.go
```

### 3. 模块内聚
相关功能放在同一模块：
- DTO 和 Service 在同一层
- Repository 接口和实现在对应层
- Handler 和 Route 在同一子模块

## 文件组织模式

### 典型的 Go 文件结构
```go
package main  // 或其他包名

import (
    // 标准库
    "context"
    "time"

    // 第三方库
    "go.uber.org/zap"

    // 内部包
    "github.com/ryuyb/fusion/internal/domain/entity"
)

// 结构体定义
type UserService struct {
    repo   repository.UserRepository
    logger *zap.Logger
}

// 构造函数
func NewUserService(
    repo repository.UserRepository,
    logger *zap.Logger,
) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

// 方法
func (s *UserService) CreateUser(ctx context.Context, dto *CreateUserDTO) error {
    // 实现
}
```

### 导入顺序
1. **标准库**: `context`, `time`, `fmt` 等
2. **第三方库**: `go.uber.org/zap`, `github.com/...`
3. **内部包**: 按模块层次组织
   - `github.com/ryuyb/fusion/internal/domain/...`
   - `github.com/ryuyb/fusion/internal/application/...`
   - `github.com/ryuyb/fusion/internal/infrastructure/...`
   - `github.com/ryuyb/fusion/internal/interface/...`
   - `github.com/ryuyb/fusion/internal/pkg/...`

## 测试组织

### 测试结构
```
test/
└── integration/
    ├── repository/             # 仓储层集成测试
    │   └── user_repository_test.go
    ├── service/                # 服务层集成测试
    │   └── user_service_test.go
    └── testutil/               # 测试工具
        └── enttest.go         # EntGO 测试工具
```

### Mock 组织
```
application/service/mocks/
└── mock_user_repository.go     # 用户仓储 Mock
```

## 配置文件组织

### 配置层次
1. **默认配置**: `configs/config.yaml`
2. **环境特定**: `configs/config.{env}.yaml`
   - `config.dev.yaml`
   - `config.prod.yaml`
3. **环境变量**: `.env` (开发)
4. **示例文件**: `.env.example`

### 配置键命名
- **使用点分隔符**: `fusion.server.port`
- **层级清晰**: `fusion.database.dsn`
- **环境敏感**: `fusion.jwt.secret`

## 文档组织

### 文档位置
- **项目根**: `CLAUDE.md`, `README.md`
- **API 文档**: `docs/api/`
- **开发指南**: `docs/development/`
- **规范文档**: `.spec-workflow/`

### 注释规范
- **导出函数**: 必需注释
  ```go
  // CreateUser creates a new user.
  func (s *UserService) CreateUser(dto *CreateUserDTO) error
  ```

- **结构体和字段**: 清晰说明
  ```go
  // User represents a user entity.
  type User struct {
      // ID is the unique identifier.
      ID int `json:"id"`
  }
  ```

## 代码大小指南

### 文件大小
- **单个文件**: < 500 行（建议）
- **特殊情况**: 可以超过，但需要良好的组织

### 函数大小
- **单个函数**: < 50 行（建议）
- **保持简单**: 一个函数只做一件事

### 结构体大小
- **字段数量**: < 20 个（建议）
- **如果过多**: 考虑拆分

### 嵌套深度
- **最大嵌套**: 3-4 层
- **避免深度嵌套**: 使用早期返回

## 包管理指南

### 包职责
- **domain**: 业务逻辑，无外部依赖
- **application**: 应用逻辑，轻依赖
- **infrastructure**: 外部依赖，重依赖
- **interface**: HTTP 逻辑
- **pkg**: 共享工具

### 依赖方向
- **上层可以依赖下层**: Interface → Application → Domain
- **下层不能依赖上层**: Domain 不能依赖 Application
- **跨层依赖通过接口**: 通过接口解耦

## 最佳实践

### DO
✅ 按层组织代码
✅ 使用 FX 进行依赖注入
✅ 依赖接口而非实现
✅ 每个模块有清晰的职责
✅ 遵循 Clean Architecture
✅ 文件命名清晰一致
✅ 测试文件靠近源码
✅ 使用 Mock 进行单元测试

### DON'T
❌ 按功能组织代码（跨层）
❌ 层间循环依赖
❌ 基础设施依赖应用层
❌ 在 Domain 层使用外部库
❌ 在接口层包含业务逻辑
❌ 文件过大或过小
❌ 缺少测试
❌ 随意导入包

---

**文档版本**: v1.0
**最后更新**: 2025-11-09
**维护者**: Fusion 架构团队