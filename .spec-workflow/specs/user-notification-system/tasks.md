# User Notification System - Implementation Tasks

本文档将用户通知系统的设计分解为可执行的原子任务。每个任务都包含具体的文件路径、需求关联、实现细节和验收标准。

---

## 阶段 1: 配置扩展 (Priority: 1)

### Task 1.1: 扩展配置结构
**文件路径**: `internal/infrastructure/config/config.go`
**关联需求**: 需求文档 2.3 - 通知偏好管理
**关联设计**: 设计文档 3.3 - 配置层次说明

#### 实现细节
在 Config 结构体中添加通知配置：

```go
type Config struct {
    // ... 现有字段 ...

    // Notification Configuration
    Notification NotificationConfig `mapstructure:"notification"`
}

type NotificationConfig struct {
    Email    EmailConfig    `mapstructure:"email"`
    Push     PushConfig     `mapstructure:"push"`
    Webhook  WebhookConfig  `mapstructure:"webhook"`
    Slack    SlackConfig    `mapstructure:"slack"`    // 可扩展
    Discord  DiscordConfig  `mapstructure:"discord"`  // 可扩展
    SMS      SMSConfig      `mapstructure:"sms"`      // 可扩展
}

type EmailConfig struct {
    SMTPHost     string `mapstructure:"smtp_host"`
    SMTPPort     int    `mapstructure:"smtp_port"`
    Username     string `mapstructure:"username"`
    Password     string `mapstructure:"password"`
    FromAddress  string `mapstructure:"from_address"`
    FromName     string `mapstructure:"from_name"`
}

type PushConfig struct {
    VAPIDPublicKey  string `mapstructure:"vapid_public_key"`
    VAPIDPrivateKey string `mapstructure:"vapid_private_key"`
    VAPIDSubject    string `mapstructure:"vapid_subject"`
}

type WebhookConfig struct {
    SignatureSecret string `mapstructure:"signature_secret"`
    TimeoutSeconds  int    `mapstructure:"timeout_seconds"`
    MaxRetries      int    `mapstructure:"max_retries"`
}

// 可扩展更多渠道配置...
```

**依赖关系**: 无
**验收标准**:
- [ ] Config 结构体包含完整的通知配置
- [ ] 支持 YAML 配置和环境变量
- [ ] 向后兼容现有配置

---

## 阶段 2: 数据模型 (Priority: 1)

### Task 2.1: 创建 Ent Schema - NotificationUserPreference
**文件路径**: `internal/infrastructure/database/schema/notification_user_preference.go`
**关联需求**: 需求文档 2.3 - 通知偏好管理
**关联设计**: 设计文档 3.1 - NotificationUserPreference 实体

#### 实现细节
定义用户通知偏好表结构：

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/orm"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/mixin"
)

// NotificationUserPreference 用户通知偏好设置
type NotificationUserPreference struct {
    ent.Schema
}

func (NotificationUserPreference) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixin.TimeMixin{}, // 包含 created_at, updated_at
        mixin.SoftDeleteMixin{},
    }
}

func (NotificationUserPreference) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id"),
        field.Int64("user_id").Unique().Comment("用户ID"),
        field.Bool("quiet_hours_enabled").Default(false).Comment("是否启用免打扰时间"),
        field.String("quiet_hours_start").Default("22:00").Comment("免打扰开始时间 HH:MM"),
        field.String("quiet_hours_end").Default("08:00").Comment("免打扰结束时间 HH:MM"),
        field.JSON("platform_filters", []string{}).Comment("平台白名单"),
        field.Int("max_notifications_per_min").Default(5).Comment("每分钟最大通知数"),
    }
}

func (NotificationUserPreference) Edges() []ent.Edge {
    return []ent.Edge{
        // 关联到用户
        edge.To("user", User.Type).Field("user_id").Required().Unique(),
    }
}

func (NotificationUserPreference) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("user_id"),
    }
}
```

**依赖关系**: Task 1.1
**验收标准**:
- [ ] 字段定义符合设计文档
- [ ] 软删除 mixin 已集成
- [ ] 索引正确设置
- [ ] 注释清晰

### Task 2.2: 创建 Ent Schema - NotificationHistory
**文件路径**: `internal/infrastructure/database/schema/notification_history.go`
**关联需求**: 需求文档 2.4 - 通知历史
**关联设计**: 设计文档 3.2 - NotificationHistory 实体

#### 实现细节
定义通知历史记录表：

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/mixin"
)

// NotificationHistory 通知历史记录
type NotificationHistory struct {
    ent.Schema
}

func (NotificationHistory) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixin.TimeMixin{},
        mixin.SoftDeleteMixin{},
    }
}

func (NotificationHistory) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id"),
        field.Int64("user_id").Comment("用户ID"),
        field.String("channel").Comment("通知渠道: email/push/webhook/slack/discord/sms"),
        field.String("status").Comment("状态: sent/delivered/read/failed"),
        field.Time("sent_at").Comment("发送时间"),
        field.Time("delivered_at").Optional().Comment("送达时间"),
        field.Time("read_at").Optional().Comment("阅读时间"),
        field.String("error_message").Optional().Comment("错误信息"),
        field.Int("retry_count").Default(0).Comment("重试次数"),
        field.JSON("metadata", map[string]interface{}{}).Comment("主播相关信息存储在此"),
    }
}

func (NotificationHistory) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).Field("user_id").Required(),
    }
}

func (NotificationHistory) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("user_id"),
        index.Fields("status"),
        index.Fields("sent_at"),
    }
}
```

**依赖关系**: Task 2.1
**验收标准**:
- [ ] Metadata 字段支持 JSON 存储
- [ ] 状态字段有索引
- [ ] 时间字段有索引

### Task 2.3: 创建 Ent Schema - NotificationChannel
**文件路径**: `internal/infrastructure/database/schema/notification_channel.go`
**关联需求**: 需求文档 2.3 - 通知渠道偏好
**关联设计**: 设计文档 3.3 - NotificationChannel 实体

#### 实现细节
定义用户级通知渠道配置表：

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/mixin"
)

// NotificationChannel 用户级通知渠道配置
type NotificationChannel struct {
    ent.Schema
}

func (NotificationChannel) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixin.TimeMixin{},
        mixin.SoftDeleteMixin{},
    }
}

func (NotificationChannel) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id"),
        field.Int64("user_id").Comment("用户ID"),
        field.String("channel").Comment("渠道名称: email/push/webhook/slack/discord/sms"),
        field.JSON("config", map[string]string{}).Comment("用户级配置存储在此"),
        field.Bool("enabled").Default(true).Comment("用户是否启用此渠道"),
        field.Int("priority").Default(5).Comment("优先级 1-10"),
    }
}

func (NotificationChannel) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).Field("user_id").Required(),
    }
}

func (NotificationChannel) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("user_id", "channel").Unique(),
    }
}
```

**依赖关系**: Task 2.2
**验收标准**:
- [ ] Config 字段支持 JSON 存储
- [ ] user_id + channel 唯一索引
- [ ] 优先级字段默认值为 5

### Task 2.4: 扩展 Subscription Schema
**文件路径**: `internal/infrastructure/database/schema/subscription.go`
**关联需求**: 需求文档 2.2 - 订阅主播开播提醒管理
**关联设计**: 设计文档 3.4 - Subscription 实体

#### 实现细节
扩展现有订阅表，添加开播提醒开关：

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/mixin"
)

// Subscription 订阅表（扩展）
type Subscription struct {
    ent.Schema
}

func (Subscription) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixin.TimeMixin{},
        mixin.SoftDeleteMixin{},
    }
}

func (Subscription) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id"),
        field.Int64("user_id").Comment("用户ID"),
        field.Int64("streamer_id").Comment("主播ID"),
        field.String("streamer_name").Comment("主播名称"),
        field.String("platform").Comment("直播平台"),
        field.Bool("notify_enabled").Default(true).Comment("开播提醒开关"),
    }
}

func (Subscription) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).Field("user_id").Required(),
    }
}

func (Subscription) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("user_id"),
        index.Fields("streamer_id"),
        index.Fields("user_id", "streamer_id").Unique(),
    }
}
```

**依赖关系**: 现有订阅表存在
**验收标准**:
- [ ] NotifyEnabled 字段默认为 true
- [ ] user_id + streamer_id 唯一索引
- [ ] 保持现有功能不变

### Task 2.5: 生成 Ent 代码
**文件路径**: 内部
**关联需求**: 无
**关联设计**: 无

#### 实现细节
运行 Ent 代码生成：

```bash
# 生成所有 schema 代码
make generate-ent

# 或者直接运行
go generate ./internal/infrastructure/database/ent
```

**依赖关系**: Task 2.1, Task 2.2, Task 2.3, Task 2.4
**验收标准**:
- [ ] 所有 .go 文件生成成功
- [ ] 没有编译错误
- [ ] 类型安全检查通过

---

## 阶段 3: Domain 层 (Priority: 2)

### Task 3.1: 创建实体 - NotificationUserPreference
**文件路径**: `internal/domain/entity/notification_user_preference.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 3.1

#### 实现细节
定义领域实体：

```go
package entity

import "time"

type NotificationUserPreference struct {
    ID                    int64
    UserID                int64
    QuietHoursEnabled     bool
    QuietHoursStart       string
    QuietHoursEnd         string
    PlatformFilters       []string
    MaxNotificationsPerMin int
    CreatedAt             time.Time
    UpdatedAt             time.Time
}
```

**依赖关系**: Task 2.1
**验收标准**:
- [ ] 字段与 schema 一致
- [ ] 无外部依赖
- [ ] 可被序列化

### Task 3.2: 创建实体 - NotificationHistory
**文件路径**: `internal/domain/entity/notification_history.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 3.2

#### 实现细节
定义通知历史实体：

```go
package entity

import "time"

type NotificationHistory struct {
    ID            int64
    UserID        int64
    Channel       string
    Status        NotificationStatus
    SentAt        time.Time
    DeliveredAt   *time.Time
    ReadAt        *time.Time
    ErrorMessage  string
    RetryCount    int
    Metadata      map[string]interface{}
    CreatedAt     time.Time
}

type NotificationStatus string

const (
    StatusSent      NotificationStatus = "sent"
    StatusDelivered NotificationStatus = "delivered"
    StatusRead      NotificationStatus = "read"
    StatusFailed    NotificationStatus = "failed"
)
```

**依赖关系**: Task 2.2
**验收标准**:
- [ ] 包含所有必要字段
- [ ] 状态类型定义清晰
- [ ] Metadata 支持任意 map

### Task 3.3: 创建实体 - NotificationChannel
**文件路径**: `internal/domain/entity/notification_channel.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 3.3

#### 实现细节
定义通知渠道实体：

```go
package entity

import "time"

type NotificationChannel struct {
    ID          int64
    UserID      int64
    Channel     string
    Config      map[string]string
    Enabled     bool
    Priority    int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**依赖关系**: Task 2.3
**验收标准**:
- [ ] Config 字段为 map[string]string
- [ ] 优先级范围 1-10

### Task 3.4: 创建实体 - Notification
**文件路径**: `internal/domain/entity/notification.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.3

#### 实现细节
定义核心通知实体：

```go
package entity

type Notification struct {
    UserID       int64
    StreamerID   int64
    StreamerName string
    Platform     string
    Title        string
    URL          string
    Thumbnail    string
    Metadata     map[string]interface{}
}
```

**依赖关系**: Task 2.4
**验收标准**:
- [ ] 包含主播和直播信息
- [ ] Metadata 支持扩展

### Task 3.5: 创建仓储接口 - NotificationUserPreference
**文件路径**: `internal/domain/repository/notification_user_preference.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.1

#### 实现细节
定义偏好仓储接口：

```go
package repository

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
)

type NotificationUserPreferenceRepository interface {
    GetByUserID(ctx context.Context, userID int64) (*entity.NotificationUserPreference, error)
    Create(ctx context.Context, pref *entity.NotificationUserPreference) error
    Update(ctx context.Context, pref *entity.NotificationUserPreference) error
    Delete(ctx context.Context, id int64) error
}
```

**依赖关系**: Task 3.1
**验收标准**:
- [ ] 包含所有 CRUD 方法
- [ ] 方法签名清晰
- [ ] 支持上下文

### Task 3.6: 创建仓储接口 - NotificationHistory
**文件路径**: `internal/domain/repository/notification_history.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 2.1

#### 实现细节
定义历史仓储接口：

```go
package repository

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/entity"
)

type NotificationHistoryRepository interface {
    Create(ctx context.Context, history *entity.NotificationHistory) error
    GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.NotificationHistory, error)
    GetByID(ctx context.Context, id int64) (*entity.NotificationHistory, error)
    UpdateStatus(ctx context.Context, id int64, status entity.NotificationStatus) error
    UpdateReadStatus(ctx context.Context, id int64) error
    GetUndelivered(ctx context.Context, userID int64) ([]*entity.NotificationHistory, error)
}
```

**依赖关系**: Task 3.2
**验收标准**:
- [ ] 支持分页查询
- [ ] 支持状态更新
- [ ] 支持未送达查询

### Task 3.7: 创建仓储接口 - NotificationChannel
**文件路径**: `internal/domain/repository/notification_channel.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.1

#### 实现细节
定义渠道仓储接口：

```go
package repository

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
)

type NotificationChannelRepository interface {
    GetByUserID(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error)
    GetByUserIDAndChannel(ctx context.Context, userID int64, channel string) (*entity.NotificationChannel, error)
    Create(ctx context.Context, channel *entity.NotificationChannel) error
    Update(ctx context.Context, channel *entity.NotificationChannel) error
    UpdateEnabled(ctx context.Context, userID int64, channel string, enabled bool) error
    Delete(ctx context.Context, id int64) error
}
```

**依赖关系**: Task 3.3
**验收标准**:
- [ ] 支持按用户查询
- [ ] 支持按用户和渠道查询
- [ ] 支持启用/禁用更新

### Task 3.8: 创建领域服务 - NotificationService
**文件路径**: `internal/domain/service/notification_service.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.2

#### 实现细节
定义通知核心服务：

```go
package service

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/domain/repository"
    "go.uber.org/zap"
)

type StreamStartEvent struct {
    StreamerID   int64
    StreamerName string
    Platform     string
    Title        string
    URL          string
    Thumbnail    string
}

type NotificationService struct {
    repo          repository.NotificationRepository
    prefRepo      repository.NotificationUserPreferenceRepository
    historyRepo   repository.NotificationHistoryRepository
    channelSvc    *ChannelService
    eventBus      EventBus
    logger        *zap.Logger
}

func NewNotificationService(
    repo repository.NotificationRepository,
    prefRepo repository.NotificationUserPreferenceRepository,
    historyRepo repository.NotificationHistoryRepository,
    channelSvc *ChannelService,
    eventBus EventBus,
    logger *zap.Logger,
) *NotificationService {
    return &NotificationService{
        repo:         repo,
        prefRepo:     prefRepo,
        historyRepo:  historyRepo,
        channelSvc:   channelSvc,
        eventBus:     eventBus,
        logger:       logger,
    }
}

func (s *NotificationService) ProcessStreamStartEvent(ctx context.Context, event StreamStartEvent) error
func (s *NotificationService) GetNotificationHistory(ctx context.Context, userID int64, limit, offset int) ([]*entity.NotificationHistory, error)
func (s *NotificationService) MarkAsRead(ctx context.Context, userID, historyID int64) error
```

**依赖关系**: Task 3.4, Task 3.5, Task 3.6, Task 3.7
**验收标准**:
- [ ] 处理开播事件逻辑
- [ ] 查询历史记录
- [ ] 标记已读功能

### Task 3.9: 创建领域服务 - NotificationPreferenceService
**文件路径**: `internal/domain/service/notification_preference_service.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.2

#### 实现细节
定义偏好服务：

```go
package service

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/domain/repository"
    "go.uber.org/zap"
)

type NotificationPreferenceService struct {
    repo   repository.NotificationUserPreferenceRepository
    logger *zap.Logger
}

func NewNotificationPreferenceService(
    repo repository.NotificationUserPreferenceRepository,
    logger *zap.Logger,
) *NotificationPreferenceService {
    return &NotificationPreferenceService{
        repo:   repo,
        logger: logger,
    }
}

func (s *NotificationPreferenceService) GetPreferences(ctx context.Context, userID int64) (*entity.NotificationUserPreference, error)
func (s *NotificationPreferenceService) UpdatePreferences(ctx context.Context, userID int64, pref *entity.NotificationUserPreference) error
func (s *NotificationPreferenceService) CheckTimeWindow(ctx context.Context, userID int64, now time.Time) (bool, error)
func (s *NotificationPreferenceService) IsChannelEnabled(ctx context.Context, userID int64, channel string) (bool, error)
```

**依赖关系**: Task 3.1, Task 3.5
**验收标准**:
- [ ] 偏好 CRUD 操作
- [ ] 时间窗口检查
- [ ] 渠道启用检查

### Task 3.10: 创建领域服务 - NotificationHistoryService
**文件路径**: `internal/domain/service/notification_history_service.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 2.2

#### 实现细节
定义历史服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/domain/repository"
    "go.uber.org/zap"
)

type NotificationHistoryService struct {
    repo   repository.NotificationHistoryRepository
    logger *zap.Logger
}

func NewNotificationHistoryService(
    repo repository.NotificationHistoryRepository,
    logger *zap.Logger,
) *NotificationHistoryService {
    return &NotificationHistoryService{
        repo:   repo,
        logger: logger,
    }
}

func (s *NotificationHistoryService) GetHistory(ctx context.Context, userID int64, limit, offset int) ([]*entity.NotificationHistory, error)
func (s *NotificationHistoryService) GetByID(ctx context.Context, id int64) (*entity.NotificationHistory, error)
func (s *NotificationHistoryService) MarkAsRead(ctx context.Context, id int64) error
func (s *NotificationHistoryService) MarkMultipleAsRead(ctx context.Context, userID int64, ids []int64) error
```

**依赖关系**: Task 3.2, Task 3.6
**验收标准**:
- [ ] 历史记录查询
- [ ] 单个标记已读
- [ ] 批量标记已读

### Task 3.11: 创建领域服务 - NotificationChannelService
**文件路径**: `internal/domain/service/notification_channel_service.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.2

#### 实现细节
定义渠道服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/domain/repository"
    "go.uber.org/zap"
)

type NotificationChannelService struct {
    repo   repository.NotificationChannelRepository
    logger *zap.Logger
}

func NewNotificationChannelService(
    repo repository.NotificationChannelRepository,
    logger *zap.Logger,
) *NotificationChannelService {
    return &NotificationChannelService{
        repo:   repo,
        logger: logger,
    }
}

func (s *NotificationChannelService) GetChannels(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error)
func (s *NotificationChannelService) GetChannel(ctx context.Context, userID int64, channel string) (*entity.NotificationChannel, error)
func (s *NotificationChannelService) UpdateChannel(ctx context.Context, channel *entity.NotificationChannel) error
func (s *NotificationChannelService) EnableChannel(ctx context.Context, userID int64, channel string) error
func (s *NotificationChannelService) DisableChannel(ctx context.Context, userID int64, channel string) error
```

**依赖关系**: Task 3.3, Task 3.7
**验收标准**:
- [ ] 渠道配置管理
- [ ] 渠道启用/禁用
- [ ] 按用户和渠道查询

### Task 3.12: 创建领域服务 - SubscriptionNotificationService
**文件路径**: `internal/domain/service/subscription_notification_service.go`
**关联需求**: 需求文档 2.2
**关联设计**: 设计文档 2.2

#### 实现细节
定义订阅通知服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "go.uber.org/zap"
)

type SubscriptionNotifyStatus struct {
    ID             int64
    StreamerID     int64
    StreamerName   string
    Platform       string
    NotifyEnabled  bool
}

type BatchUpdateRequest struct {
    SubscriptionIDs []int64
    NotifyEnabled   bool
}

type SubscriptionNotificationService struct {
    notifSvc *NotificationService
    prefSvc  *NotificationPreferenceService
    repo     SubscriptionRepository
    logger   *zap.Logger
}

func NewSubscriptionNotificationService(
    notifSvc *NotificationService,
    prefSvc *NotificationPreferenceService,
    repo SubscriptionRepository,
    logger *zap.Logger,
) *SubscriptionNotificationService {
    return &SubscriptionNotificationService{
        notifSvc: notifSvc,
        prefSvc:  prefSvc,
        repo:     repo,
        logger:   logger,
    }
}

func (s *SubscriptionNotificationService) UpdateStreamNotifyStatus(ctx context.Context, userID, subscriptionID int64, enabled bool) error
func (s *SubscriptionNotificationService) BatchUpdateNotifyStatus(ctx context.Context, userID int64, req BatchUpdateRequest) error
func (s *SubscriptionNotificationService) GetSubscriptionNotifyStatus(ctx context.Context, userID int64) ([]SubscriptionNotifyStatus, error)
```

**依赖关系**: Task 3.8, Task 3.9
**验收标准**:
- [ ] 单个提醒状态更新
- [ ] 批量提醒状态更新
- [ ] 查询提醒状态列表

### Task 3.13: 创建领域服务 - NotificationChannel
**文件路径**: `internal/domain/service/notification_channel.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.4

#### 实现细节
定义基础发送器接口和渠道管理：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
)

type NotificationRequest struct {
    UserID       int64
    StreamerID   int64
    StreamerName string
    Platform     string
    Title        string
    URL          string
    Thumbnail    string
    Metadata     map[string]interface{}
}

type ChannelInfo struct {
    Name        string
    DisplayName string
    Description string
    Config      map[string]string
    Enabled     bool
}

type BaseSender interface {
    Channel() string
    Send(ctx context.Context, req NotificationRequest) error
    Validate(req NotificationRequest) error
    GetChannelInfo() ChannelInfo
}

type ChannelService struct {
    registry *ChannelRegistry
    logger   *zap.Logger
}

func NewChannelService(registry *ChannelRegistry, logger *zap.Logger) *ChannelService {
    return &ChannelService{
        registry: registry,
        logger:   logger,
    }
}

func (s *ChannelService) RegisterChannel(sender BaseSender) error
func (s *ChannelService) GetChannel(channel string) (BaseSender, bool)
func (s *ChannelService) ListChannels() []ChannelInfo
func (s *ChannelService) UpdateChannelConfig(channel string, config map[string]string) error
func (s *ChannelService) EnableChannel(channel string) error
func (s *ChannelService) DisableChannel(channel string) error
```

**依赖关系**: Task 2.4
**验收标准**:
- [ ] BaseSender 接口定义
- [ ] ChannelService 功能完整
- [ ] 支持渠道注册和管理

### Task 3.14: 创建领域模块
**文件路径**: `internal/domain/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
导出领域层 FX 模块：

```go
package domain

import "go.uber.org/fx"

var Module = fx.Module(
    "domain",
    // 领域层依赖
)
```

**依赖关系**: Task 3.1-3.13
**验收标准**:
- [ ] Module 正确导出
- [ ] 可被应用层依赖

---

## 阶段 4: Application 层 (Priority: 2)

### Task 4.1: 创建 DTO - 请求
**文件路径**: `internal/application/dto/request/notification.go`
**关联需求**: 需求文档 2.2, 2.3
**关联设计**: 设计文档 2.5

#### 实现细节
定义请求 DTO：

```go
package request

import "time"

type UpdatePreferencesRequest struct {
    QuietHoursEnabled     bool     `json:"quiet_hours_enabled" validate:"required"`
    QuietHoursStart       string   `json:"quiet_hours_start" validate:"required,len=5"`
    QuietHoursEnd         string   `json:"quiet_hours_end" validate:"required,len=5"`
    PlatformFilters       []string `json:"platform_filters"`
    MaxNotificationsPerMin int     `json:"max_notifications_per_min" validate:"required,min=1,max=60"`
}

type UpdateChannelConfigRequest struct {
    Config   map[string]string `json:"config" validate:"required"`
    Enabled  *bool            `json:"enabled"`
    Priority *int             `json:"priority" validate:"min=1,max=10"`
}

type BatchUpdateNotifyRequest struct {
    SubscriptionIDs []int64 `json:"subscription_ids" validate:"required,min=1"`
    NotifyEnabled   bool    `json:"notify_enabled" validate:"required"`
}
```

**依赖关系**: Task 3.9, Task 3.11, Task 3.12
**验收标准**:
- [ ] 字段验证标签正确
- [ ] JSON 标签匹配
- [ ] 验证规则完整

### Task 4.2: 创建 DTO - 响应
**文件路径**: `internal/application/dto/response/notification.go`
**关联需求**: 需求文档 2.2, 2.3, 2.4
**关联设计**: 设计文档 2.5

#### 实现细节
定义响应 DTO：

```go
package response

import (
    "time"
    "github.com/ryuyb/fusion/internal/domain/entity"
)

type NotificationHistoryItem struct {
    ID            int64                  `json:"id"`
    Channel       string                 `json:"channel"`
    Status        string                 `json:"status"`
    SentAt        time.Time              `json:"sent_at"`
    DeliveredAt   *time.Time             `json:"delivered_at,omitempty"`
    ReadAt        *time.Time             `json:"read_at,omitempty"`
    ErrorMessage  string                 `json:"error_message,omitempty"`
    RetryCount    int                    `json:"retry_count"`
    Metadata      map[string]interface{} `json:"metadata"`
}

type NotificationPreferences struct {
    ID                    int64     `json:"id"`
    UserID                int64     `json:"user_id"`
    QuietHoursEnabled     bool      `json:"quiet_hours_enabled"`
    QuietHoursStart       string    `json:"quiet_hours_start"`
    QuietHoursEnd         string    `json:"quiet_hours_end"`
    PlatformFilters       []string  `json:"platform_filters"`
    MaxNotificationsPerMin int      `json:"max_notifications_per_min"`
    CreatedAt             time.Time `json:"created_at"`
    UpdatedAt             time.Time `json:"updated_at"`
}

type NotificationChannelItem struct {
    ID          int64             `json:"id"`
    Channel     string            `json:"channel"`
    Config      map[string]string `json:"config"`
    Enabled     bool              `json:"enabled"`
    Priority    int               `json:"priority"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

type SubscriptionNotifyStatusItem struct {
    ID            int64  `json:"id"`
    StreamerID    int64  `json:"streamer_id"`
    StreamerName  string `json:"streamer_name"`
    Platform      string `json:"platform"`
    NotifyEnabled bool   `json:"notify_enabled"`
}
```

**依赖关系**: Task 3.2, Task 3.1, Task 3.3, Task 3.12
**验收标准**:
- [ ] 字段与实体一致
- [ ] JSON 序列化正确
- [ ] 包含必要字段

### Task 4.3: 创建应用服务 - NotificationService
**文件路径**: `internal/application/service/notification_service.go`
**关联需求**: 需求文档 2.1, 2.4
**关联设计**: 设计文档 2.6

#### 实现细节
定义应用层通知服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/dto/response"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/domain/repository"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
)

type NotificationService struct {
    notifSvc  *service.NotificationService
    historySvc *service.NotificationHistoryService
    prefSvc   *service.NotificationPreferenceService
    logger    *zap.Logger
}

func NewNotificationService(
    notifSvc *service.NotificationService,
    historySvc *service.NotificationHistoryService,
    prefSvc *service.NotificationPreferenceService,
    logger *zap.Logger,
) *NotificationService {
    return &NotificationService{
        notifSvc:  notifSvc,
        historySvc: historySvc,
        prefSvc:   prefSvc,
        logger:    logger,
    }
}

func (s *NotificationService) GetHistory(ctx context.Context, userID int64, limit, offset int) ([]response.NotificationHistoryItem, error)
func (s *NotificationService) MarkAsRead(ctx context.Context, userID, historyID int64) error
func (s *NotificationService) MarkMultipleAsRead(ctx context.Context, userID int64, ids []int64) error
```

**依赖关系**: Task 3.8, Task 3.10, Task 3.9
**验收标准**:
- [ ] 领域服务封装
- [ ] DTO 转换
- [ ] 错误处理

### Task 4.4: 创建应用服务 - PreferenceService
**文件路径**: `internal/application/service/notification_preference_service.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.6

#### 实现细节
定义应用层偏好服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/dto/response"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
)

type NotificationPreferenceService struct {
    prefSvc *service.NotificationPreferenceService
    logger  *zap.Logger
}

func NewNotificationPreferenceService(
    prefSvc *service.NotificationPreferenceService,
    logger *zap.Logger,
) *NotificationPreferenceService {
    return &NotificationPreferenceService{
        prefSvc: prefSvc,
        logger:  logger,
    }
}

func (s *NotificationPreferenceService) GetPreferences(ctx context.Context, userID int64) (*response.NotificationPreferences, error)
func (s *NotificationPreferenceService) UpdatePreferences(ctx context.Context, userID int64, req *request.UpdatePreferencesRequest) error
```

**依赖关系**: Task 3.9
**验收标准**:
- [ ] 领域服务封装
- [ ] DTO 转换和验证
- [ ] 业务逻辑调用

### Task 4.5: 创建应用服务 - ChannelService
**文件路径**: `internal/application/service/notification_channel_service.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.6

#### 实现细节
定义应用层渠道服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/dto/response"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
)

type NotificationChannelService struct {
    channelSvc *service.ChannelService
    logger     *zap.Logger
}

func NewNotificationChannelService(
    channelSvc *service.ChannelService,
    logger *zap.Logger,
) *NotificationChannelService {
    return &NotificationChannelService{
        channelSvc: channelSvc,
        logger:     logger,
    }
}

func (s *NotificationChannelService) GetChannels(ctx context.Context, userID int64) ([]response.NotificationChannelItem, error)
func (s *NotificationChannelService) UpdateChannel(ctx context.Context, userID int64, channel string, req *request.UpdateChannelConfigRequest) error
func (s *NotificationChannelService) EnableChannel(ctx context.Context, userID int64, channel string) error
func (s *NotificationChannelService) DisableChannel(ctx context.Context, userID int64, channel string) error
```

**依赖关系**: Task 3.11
**验收标准**:
- [ ] 领域服务封装
- [ ] 用户权限检查
- [ ] 业务逻辑调用

### Task 4.6: 创建应用服务 - SubscriptionNotificationService
**文件路径**: `internal/application/service/subscription_notification_service.go`
**关联需求**: 需求文档 2.2
**关联设计**: 设计文档 2.6

#### 实现细节
定义应用层订阅通知服务：

```go
package service

import (
    "context"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/dto/response"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
)

type SubscriptionNotificationService struct {
    notifSvc *service.SubscriptionNotificationService
    logger   *zap.Logger
}

func NewSubscriptionNotificationService(
    notifSvc *service.SubscriptionNotificationService,
    logger *zap.Logger,
) *SubscriptionNotificationService {
    return &SubscriptionNotificationService{
        notifSvc: notifSvc,
        logger:   logger,
    }
}

func (s *SubscriptionNotificationService) UpdateStreamNotifyStatus(ctx context.Context, userID, subscriptionID int64, enabled bool) error
func (s *SubscriptionNotificationService) BatchUpdateNotifyStatus(ctx context.Context, userID int64, req *request.BatchUpdateNotifyRequest) error
func (s *SubscriptionNotificationService) GetSubscriptionNotifyStatus(ctx context.Context, userID int64) ([]response.SubscriptionNotifyStatusItem, error)
```

**依赖关系**: Task 3.12
**验收标准**:
- [ ] 领域服务封装
- [ ] 批量操作支持
- [ ] 事务保证

### Task 4.7: 创建应用模块
**文件路径**: `internal/application/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
导出应用层 FX 模块：

```go
package application

import "go.uber.org/fx"

var Module = fx.Module(
    "application",
    // 应用层依赖
)
```

**依赖关系**: Task 4.1-4.6
**验收标准**:
- [ ] Module 正确导出
- [ ] 可被接口层依赖

---

## 阶段 5: Infrastructure 层 (Priority: 3)

### Task 5.1: 创建仓储实现 - NotificationUserPreference
**文件路径**: `internal/infrastructure/repository/notification_user_preference_repository.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.7

#### 实现细节
实现偏好仓储：

```go
package repository

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationuserpreference"
)

type NotificationUserPreferenceRepository struct {
    client *ent.Client
}

func NewNotificationUserPreferenceRepository(client *ent.Client) *NotificationUserPreferenceRepository {
    return &NotificationUserPreferenceRepository{
        client: client,
    }
}

func (r *NotificationUserPreferenceRepository) GetByUserID(ctx context.Context, userID int64) (*entity.NotificationUserPreference, error)
func (r *NotificationUserPreferenceRepository) Create(ctx context.Context, pref *entity.NotificationUserPreference) error
func (r *NotificationUserPreferenceRepository) Update(ctx context.Context, pref *entity.NotificationUserPreference) error
func (r *NotificationUserPreferenceRepository) Delete(ctx context.Context, id int64) error
```

**依赖关系**: Task 3.5
**验收标准**:
- [ ] Ent 客户端使用正确
- [ ] 实体转换正确
- [ ] 错误处理完整

### Task 5.2: 创建仓储实现 - NotificationHistory
**文件路径**: `internal/infrastructure/repository/notification_history_repository.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 2.7

#### 实现细节
实现历史仓储：

```go
package repository

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationhistory"
)

type NotificationHistoryRepository struct {
    client *ent.Client
}

func NewNotificationHistoryRepository(client *ent.Client) *NotificationHistoryRepository {
    return &NotificationHistoryRepository{
        client: client,
    }
}

func (r *NotificationHistoryRepository) Create(ctx context.Context, history *entity.NotificationHistory) error
func (r *NotificationHistoryRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.NotificationHistory, error)
func (r *NotificationHistoryRepository) GetByID(ctx context.Context, id int64) (*entity.NotificationHistory, error)
func (r *NotificationHistoryRepository) UpdateStatus(ctx context.Context, id int64, status entity.NotificationStatus) error
func (r *NotificationHistoryRepository) UpdateReadStatus(ctx context.Context, id int64) error
func (r *NotificationHistoryRepository) GetUndelivered(ctx context.Context, userID int64) ([]*entity.NotificationHistory, error)
```

**依赖关系**: Task 3.6
**验收标准**:
- [ ] 分页查询正确
- [ ] 状态更新事务性
- [ ] 索引使用优化

### Task 5.3: 创建仓储实现 - NotificationChannel
**文件路径**: `internal/infrastructure/repository/notification_channel_repository.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.7

#### 实现细节
实现渠道仓储：

```go
package repository

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/entity"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationchannel"
)

type NotificationChannelRepository struct {
    client *ent.Client
}

func NewNotificationChannelRepository(client *ent.Client) *NotificationChannelRepository {
    return &NotificationChannelRepository{
        client: client,
    }
}

func (r *NotificationChannelRepository) GetByUserID(ctx context.Context, userID int64) ([]*entity.NotificationChannel, error)
func (r *NotificationChannelRepository) GetByUserIDAndChannel(ctx context.Context, userID int64, channel string) (*entity.NotificationChannel, error)
func (r *NotificationChannelRepository) Create(ctx context.Context, channel *entity.NotificationChannel) error
func (r *NotificationChannelRepository) Update(ctx context.Context, channel *entity.NotificationChannel) error
func (r *NotificationChannelRepository) UpdateEnabled(ctx context.Context, userID int64, channel string, enabled bool) error
func (r *NotificationChannelRepository) Delete(ctx context.Context, id int64) error
```

**依赖关系**: Task 3.7
**验收标准**:
- [ ] 唯一索引使用
- [ ] 条件更新优化
- [ ] 事务保证

### Task 5.4: 创建基础发送器接口
**文件路径**: `internal/infrastructure/sender/base_sender.go`
**关联需求**: 需求文档 2.1, 2.3
**关联设计**: 设计文档 2.4

#### 实现细节
定义基础发送器：

```go
package sender

import (
    "context"
    "github.com/ryuyb/fusion/internal/domain/service"
)

type BaseSender interface {
    Channel() string
    Send(ctx context.Context, req service.NotificationRequest) error
    Validate(req service.NotificationRequest) error
    GetChannelInfo() service.ChannelInfo
}

type ChannelRegistry struct {
    senders map[string]BaseSender
    logger  *zap.Logger
    mutex   sync.RWMutex
}

func NewChannelRegistry(logger *zap.Logger) *ChannelRegistry {
    return &ChannelRegistry{
        senders: make(map[string]BaseSender),
        logger:  logger,
    }
}

func (r *ChannelRegistry) Register(sender BaseSender) error
func (r *ChannelRegistry) GetSender(channel string) (BaseSender, bool)
func (r *ChannelRegistry) ListChannels() []string
func (r *ChannelRegistry) GetConfig(channel string) (*ChannelConfig, bool)
func (r *ChannelRegistry) UpdateConfig(channel string, config *ChannelConfig) error
```

**依赖关系**: Task 3.13
**验收标准**:
- [ ] 线程安全
- [ ] 接口定义完整
- [ ] 注册机制正确

### Task 5.5: 创建发送器 - Email
**文件路径**: `internal/infrastructure/sender/email_sender.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.4

#### 实现细节
实现邮件发送器：

```go
package sender

import (
    "context"
    "fmt"
    "time"
    "github.com/ryuyb/fusion/internal/domain/service"
    "github.com/ryuyb/fusion/internal/infrastructure/config"
    "go.uber.org/zap"
    "gopkg.in/gomail.v2"
)

type EmailSender struct {
    smtpClient *gomail.Dialer
    template   *EmailTemplate
    logger     *zap.Logger
}

func NewEmailSender(cfg config.EmailConfig, logger *zap.Logger) *EmailSender {
    dialer := gomail.NewDialer(
        cfg.SMTPHost,
        cfg.SMTPPort,
        cfg.Username,
        cfg.Password,
    )

    return &EmailSender{
        smtpClient: dialer,
        template:   NewEmailTemplate(),
        logger:     logger,
    }
}

func (s *EmailSender) Channel() string {
    return "email"
}

func (s *EmailSender) Send(ctx context.Context, req service.NotificationRequest) error
func (s *EmailSender) Validate(req service.NotificationRequest) error
func (s *EmailSender) GetChannelInfo() service.ChannelInfo
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] SMTP 连接正确
- [ ] 邮件模板渲染
- [ ] 错误处理完整

### Task 5.6: 创建发送器 - Push
**文件路径**: `internal/infrastructure/sender/push_sender.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.4

#### 实现细节
实现推送发送器：

```go
package sender

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/service"
    "github.com/ryuyb/fusion/internal/infrastructure/config"
    "go.uber.org/zap"
    "github.com/segmentio/ksuid"
)

type PushSender struct {
    vapidPublicKey  string
    vapidPrivateKey string
    httpClient      *resty.Client
    logger          *zap.Logger
}

func NewPushSender(cfg config.PushConfig, logger *zap.Logger) *PushSender {
    return &PushSender{
        vapidPublicKey:  cfg.VAPIDPublicKey,
        vapidPrivateKey: cfg.VAPIDPrivateKey,
        httpClient:      resty.New(),
        logger:          logger,
    }
}

func (s *PushSender) Channel() string {
    return "push"
}

func (s *PushSender) Send(ctx context.Context, req service.NotificationRequest) error
func (s *PushSender) Validate(req service.NotificationRequest) error
func (s *PushSender) GetChannelInfo() service.ChannelInfo
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] VAPID 密钥配置
- [ ] 推送载荷构造
- [ ] 订阅管理

### Task 5.7: 创建发送器 - Webhook
**文件路径**: `internal/infrastructure/sender/webhook_sender.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.4

#### 实现细节
实现 Webhook 发送器：

```go
package sender

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "time"
    "github.com/ryuyb/fusion/internal/domain/service"
    "github.com/ryuyb/fusion/internal/infrastructure/config"
    "go.uber.org/zap"
    "resty.dev/v3"
)

type WebhookSender struct {
    httpClient   *resty.Client
    signature    SignatureVerifier
    retryPolicy  RetryPolicy
    logger       *zap.Logger
}

func NewWebhookSender(cfg config.WebhookConfig, logger *zap.Logger) *WebhookSender {
    return &WebhookSender{
        httpClient:   resty.New(),
        signature:    NewSignatureVerifier(cfg.SignatureSecret),
        retryPolicy:  NewRetryPolicy(cfg.MaxRetries),
        logger:       logger,
    }
}

func (s *WebhookSender) Channel() string {
    return "webhook"
}

func (s *WebhookSender) Send(ctx context.Context, req service.NotificationRequest) error
func (s *WebhookSender) Validate(req service.NotificationRequest) error
func (s *WebhookSender) GetChannelInfo() service.ChannelInfo
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] 签名验证实现
- [ ] 重试机制正确
- [ ] 错误处理完整

### Task 5.8: 创建发送器 - Slack
**文件路径**: `internal/infrastructure/sender/slack_sender.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.4

#### 实现细节
实现 Slack 发送器：

```go
package sender

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
    "resty.dev/v3"
)

type SlackSender struct {
    httpClient   *resty.Client
    botToken     string
    defaultChannel string
    logger       *zap.Logger
}

func NewSlackSender(botToken, defaultChannel string, logger *zap.Logger) *SlackSender {
    return &SlackSender{
        httpClient:   resty.New(),
        botToken:     botToken,
        defaultChannel: defaultChannel,
        logger:       logger,
    }
}

func (s *SlackSender) Channel() string {
    return "slack"
}

func (s *SlackSender) Send(ctx context.Context, req service.NotificationRequest) error
func (s *SlackSender) Validate(req service.NotificationRequest) error
func (s *SlackSender) GetChannelInfo() service.ChannelInfo
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] Slack API 调用
- [ ] 消息格式正确
- [ ] 错误处理完整

### Task 5.9: 创建发送器 - Discord
**文件路径**: `internal/infrastructure/sender/discord_sender.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.4

#### 实现细节
实现 Discord 发送器：

```go
package sender

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
    "resty.dev/v3"
)

type DiscordSender struct {
    httpClient *resty.Client
    webhookURL string
    logger     *zap.Logger
}

func NewDiscordSender(webhookURL string, logger *zap.Logger) *DiscordSender {
    return &DiscordSender{
        httpClient: resty.New(),
        webhookURL: webhookURL,
        logger:     logger,
    }
}

func (d *DiscordSender) Channel() string {
    return "discord"
}

func (d *DiscordSender) Send(ctx context.Context, req service.NotificationRequest) error
func (d *DiscordSender) Validate(req service.NotificationRequest) error
func (d *DiscordSender) GetChannelInfo() service.ChannelInfo
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] Discord Webhook 调用
- [ ] Embed 格式正确
- [ ] 错误处理完整

### Task 5.10: 创建发送器 - SMS
**文件路径**: `internal/infrastructure/sender/sms_sender.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.4

#### 实现细节
实现短信发送器：

```go
package sender

import (
    "context"
    "time"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
)

type SMSClient interface {
    Send(ctx context.Context, phoneNumber, message string) error
}

type SMSsender struct {
    smsClient SMSClient
    logger    *zap.Logger
}

func NewSMSsender(smsClient SMSClient, logger *zap.Logger) *SMSsender {
    return &SMSsender{
        smsClient: smsClient,
        logger:    logger,
    }
}

func (s *SMSsender) Channel() string {
    return "sms"
}

func (s *SMSsender) Send(ctx context.Context, req service.NotificationRequest) error
func (s *SMSsender) Validate(req service.NotificationRequest) error
func (s *SMSsender) GetChannelInfo() service.ChannelInfo
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] 短信接口调用
- [ ] 字数限制处理
- [ ] 错误处理完整

### Task 5.11: 创建渠道注册
**文件路径**: `internal/infrastructure/sender/registry.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.4

#### 实现细节
注册所有发送器：

```go
package sender

import (
    "context"
    "go.uber.org/fx"
    "github.com/ryuyb/fusion/internal/infrastructure/config"
)

func ProvideEmailSender(cfg config.EmailConfig, logger *zap.Logger) BaseSender {
    return NewEmailSender(cfg, logger)
}

func ProvidePushSender(cfg config.PushConfig, logger *zap.Logger) BaseSender {
    return NewPushSender(cfg, logger)
}

func ProvideWebhookSender(cfg config.WebhookConfig, logger *zap.Logger) BaseSender {
    return NewWebhookSender(cfg, logger)
}

// 提供其他发送器...

func RegisterChannels(registry *ChannelRegistry, senders ...BaseSender) {
    for _, sender := range senders {
        if err := registry.Register(sender); err != nil {
            panic(err)
        }
    }
}
```

**依赖关系**: Task 5.5-5.10
**验收标准**:
- [ ] 所有发送器注册
- [ ] FX 依赖注入配置
- [ ] 错误处理正确

### Task 5.12: 创建消息队列
**文件路径**: `internal/infrastructure/queue/notification_queue.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.5

#### 实现细节
实现通知队列：

```go
package queue

import (
    "context"
    "github.com/go-redis/redis/v8"
    "github.com/ryuyb/fusion/internal/domain/service"
    "go.uber.org/zap"
)

type NotificationQueue struct {
    redisClient redis.Client
    logger      *zap.Logger
}

func NewNotificationQueue(redisClient redis.Client, logger *zap.Logger) *NotificationQueue {
    return &NotificationQueue{
        redisClient: redisClient,
        logger:      logger,
    }
}

func (q *NotificationQueue) Enqueue(ctx context.Context, req service.NotificationRequest) error
func (q *NotificationQueue) Dequeue(ctx context.Context) (*service.NotificationRequest, error)
func (q *NotificationQueue) Process(ctx context.Context, handler func(req service.NotificationRequest) error) error
```

**依赖关系**: Task 3.13
**验收标准**:
- [ ] Redis 连接正确
- [ ] 队列操作原子性
- [ ] 错误处理完整

### Task 5.13: 创建基础设施模块
**文件路径**: `internal/infrastructure/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
导出基础设施层 FX 模块：

```go
package infrastructure

import "go.uber.org/fx"

var Module = fx.Module(
    "infrastructure",
    // 基础设施层依赖
)
```

**依赖关系**: Task 5.1-5.12
**验收标准**:
- [ ] Module 正确导出
- [ ] 可被应用层依赖

---

## 阶段 6: Interface 层 (Priority: 3)

### Task 6.1: 创建处理器 - Notification
**文件路径**: `internal/interface/http/handler/notification.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 2.6

#### 实现细节
实现通知历史处理器：

```go
package handler

import (
    "net/http"
    "strconv"
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/application/dto/response"
    "github.com/ryuyb/fusion/internal/application/service"
    "go.uber.org/zap"
)

type NotificationHandler struct {
    notifSvc *service.NotificationService
    logger   *zap.Logger
}

func NewNotificationHandler(
    notifSvc *service.NotificationService,
    logger *zap.Logger,
) *NotificationHandler {
    return &NotificationHandler{
        notifSvc: notifSvc,
        logger:   logger,
    }
}

// GetHistory 获取通知历史
// @Summary 获取通知历史
// @Description 获取用户的通知历史记录
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.NotificationHistoryItem
// @Router /api/notifications/history [get]
func (h *NotificationHandler) GetHistory(c fiber.Ctx) error {
    userID := c.Locals("user_id").(int64)

    limit, err := strconv.Atoi(c.Query("limit", "20"))
    if err != nil || limit <= 0 {
        limit = 20
    }

    offset, err := strconv.Atoi(c.Query("offset", "0"))
    if err != nil || offset < 0 {
        offset = 0
    }

    histories, err := h.notifSvc.GetHistory(c.Context(), userID, limit, offset)
    if err != nil {
        return err
    }

    return c.JSON(histories)
}

// MarkAsRead 标记通知为已读
func (h *NotificationHandler) MarkAsRead(c fiber.Ctx) error
// MarkMultipleAsRead 批量标记为已读
func (h *NotificationHandler) MarkMultipleAsRead(c fiber.Ctx) error
```

**依赖关系**: Task 4.3
**验收标准**:
- [ ] Swagger 文档完整
- [ ] 错误处理正确
- [ ] 认证中间件

### Task 6.2: 创建处理器 - Preference
**文件路径**: `internal/interface/http/handler/notification_preference.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.6

#### 实现细节
实现偏好处理器：

```go
package handler

import (
    "net/http"
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/service"
    "go.uber.org/zap"
)

type NotificationPreferenceHandler struct {
    prefSvc *service.NotificationPreferenceService
    logger  *zap.Logger
}

func NewNotificationPreferenceHandler(
    prefSvc *service.NotificationPreferenceService,
    logger *zap.Logger,
) *NotificationPreferenceHandler {
    return &NotificationPreferenceHandler{
        prefSvc: prefSvc,
        logger:  logger,
    }
}

// GetPreferences 获取用户偏好
func (h *NotificationPreferenceHandler) GetPreferences(c fiber.Ctx) error {
    userID := c.Locals("user_id").(int64)

    preferences, err := h.prefSvc.GetPreferences(c.Context(), userID)
    if err != nil {
        return err
    }

    return c.JSON(preferences)
}

// UpdatePreferences 更新用户偏好
func (h *NotificationPreferenceHandler) UpdatePreferences(c fiber.Ctx) error
```

**依赖关系**: Task 4.4
**验收标准**:
- [ ] 参数验证
- [ ] 业务逻辑调用
- [ ] 响应格式正确

### Task 6.3: 创建处理器 - Channel
**文件路径**: `internal/interface/http/handler/notification_channel.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.6

#### 实现细节
实现渠道处理器：

```go
package handler

import (
    "net/http"
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/service"
    "go.uber.org/zap"
)

type NotificationChannelHandler struct {
    channelSvc *service.NotificationChannelService
    logger     *zap.Logger
}

func NewNotificationChannelHandler(
    channelSvc *service.NotificationChannelService,
    logger *zap.Logger,
) *NotificationChannelHandler {
    return &NotificationChannelHandler{
        channelSvc: channelSvc,
        logger:     logger,
    }
}

// GetChannels 获取用户渠道配置
func (h *NotificationChannelHandler) GetChannels(c fiber.Ctx) error {
    userID := c.Locals("user_id").(int64)

    channels, err := h.channelSvc.GetChannels(c.Context(), userID)
    if err != nil {
        return err
    }

    return c.JSON(channels)
}

// UpdateChannelConfig 更新渠道配置
func (h *NotificationChannelHandler) UpdateChannelConfig(c fiber.Ctx) error

// EnableChannel 启用渠道
func (h *NotificationChannelHandler) EnableChannel(c fiber.Ctx) error

// DisableChannel 禁用渠道
func (h *NotificationChannelHandler) DisableChannel(c fiber.Ctx) error
```

**依赖关系**: Task 4.5
**验收标准**:
- [ ] 渠道配置管理
- [ ] 启用/禁用功能
- [ ] 用户权限检查

### Task 6.4: 创建处理器 - Subscription
**文件路径**: `internal/interface/http/handler/subscription_notification.go`
**关联需求**: 需求文档 2.2
**关联设计**: 设计文档 2.6

#### 实现细节
实现订阅通知处理器：

```go
package handler

import (
    "net/http"
    "strconv"
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/application/dto/request"
    "github.com/ryuyb/fusion/internal/application/service"
    "go.uber.org/zap"
)

type SubscriptionNotificationHandler struct {
    subSvc *service.SubscriptionNotificationService
    logger *zap.Logger
}

func NewSubscriptionNotificationHandler(
    subSvc *service.SubscriptionNotificationService,
    logger *zap.Logger,
) *SubscriptionNotificationHandler {
    return &SubscriptionNotificationHandler{
        subSvc: subSvc,
        logger: logger,
    }
}

// UpdateStreamNotifyStatus 更新主播提醒状态
func (h *SubscriptionNotificationHandler) UpdateStreamNotifyStatus(c fiber.Ctx) error {
    userID := c.Locals("user_id").(int64)
    subscriptionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "invalid subscription id")
    }

    enabled := c.Query("enabled") == "true"
    if err := h.subSvc.UpdateStreamNotifyStatus(c.Context(), userID, subscriptionID, enabled); err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "success": true,
    })
}

// BatchUpdateNotifyStatus 批量更新提醒状态
func (h *SubscriptionNotificationHandler) BatchUpdateNotifyStatus(c fiber.Ctx) error

// GetSubscriptionNotifyStatus 获取订阅提醒状态
func (h *SubscriptionNotificationHandler) GetSubscriptionNotifyStatus(c fiber.Ctx) error
```

**依赖关系**: Task 4.6
**验收标准**:
- [ ] 单个更新功能
- [ ] 批量更新功能
- [ ] 状态查询功能

### Task 6.5: 创建路由 - Notification
**文件路径**: `internal/interface/http/route/notification.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 2.6

#### 实现细节
定义通知路由：

```go
package route

import (
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/interface/http/handler"
    "github.com/ryuyb/fusion/internal/interface/http/middleware"
)

type NotificationRoute struct {
    handler *handler.NotificationHandler
}

func NewNotificationRoute(
    handler *handler.NotificationHandler,
) *NotificationRoute {
    return &NotificationRoute{
        handler: handler,
    }
}

func (r *NotificationRoute) Register(app fiber.Router) {
    group := app.Group("/api/notifications")
    group.Use(middleware.Auth()) // JWT 认证

    group.Get("/history", r.handler.GetHistory)
    group.Post("/history/:id/read", r.handler.MarkAsRead)
    group.Post("/history/batch-read", r.handler.MarkMultipleAsRead)
}
```

**依赖关系**: Task 6.1
**验收标准**:
- [ ] 路由定义正确
- [ ] 中间件配置
- [ ] 认证保护

### Task 6.6: 创建路由 - Preference
**文件路径**: `internal/interface/http/route/notification_preference.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.6

#### 实现细节
定义偏好路由：

```go
package route

import (
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/interface/http/handler"
    "github.com/ryuyb/fusion/internal/interface/http/middleware"
)

type NotificationPreferenceRoute struct {
    handler *handler.NotificationPreferenceHandler
}

func NewNotificationPreferenceRoute(
    handler *handler.NotificationPreferenceHandler,
) *NotificationPreferenceRoute {
    return &NotificationPreferenceRoute{
        handler: handler,
    }
}

func (r *NotificationPreferenceRoute) Register(app fiber.Router) {
    group := app.Group("/api/notifications/preferences")
    group.Use(middleware.Auth())

    group.Get("/", r.handler.GetPreferences)
    group.Put("/", r.handler.UpdatePreferences)
}
```

**依赖关系**: Task 6.2
**验收标准**:
- [ ] GET/PUT 路由
- [ ] 认证保护
- [ ] 验证中间件

### Task 6.7: 创建路由 - Channel
**文件路径**: `internal/interface/http/route/notification_channel.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.6

#### 实现细节
定义渠道路由：

```go
package route

import (
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/interface/http/handler"
    "github.com/ryuyb/fusion/internal/interface/http/middleware"
)

type NotificationChannelRoute struct {
    handler *handler.NotificationChannelHandler
}

func NewNotificationChannelRoute(
    handler *handler.NotificationChannelHandler,
) *NotificationChannelRoute {
    return &NotificationChannelRoute{
        handler: handler,
    }
}

func (r *NotificationChannelRoute) Register(app fiber.Router) {
    group := app.Group("/api/notifications/channels")
    group.Use(middleware.Auth())

    group.Get("/", r.handler.GetChannels)
    group.Put("/:channel", r.handler.UpdateChannelConfig)
    group.Post("/:channel/enable", r.handler.EnableChannel)
    group.Post("/:channel/disable", r.handler.DisableChannel)
}
```

**依赖关系**: Task 6.3
**验收标准**:
- [ ] 渠道 CRUD 路由
- [ ] 启用/禁用路由
- [ ] 认证保护

### Task 6.8: 创建路由 - Subscription
**文件路径**: `internal/interface/http/route/subscription_notification.go`
**关联需求**: 需求文档 2.2
**关联设计**: 设计文档 2.6

#### 实现细节
定义订阅通知路由：

```go
package route

import (
    "github.com/gofiber/fiber/v3"
    "github.com/ryuyb/fusion/internal/interface/http/handler"
    "github.com/ryuyb/fusion/internal/interface/http/middleware"
)

type SubscriptionNotificationRoute struct {
    handler *handler.SubscriptionNotificationHandler
}

func NewSubscriptionNotificationRoute(
    handler *handler.SubscriptionNotificationHandler,
) *SubscriptionNotificationRoute {
    return &SubscriptionNotificationRoute{
        handler: handler,
    }
}

func (r *SubscriptionNotificationRoute) Register(app fiber.Router) {
    group := app.Group("/api/subscriptions")
    group.Use(middleware.Auth())

    group.Get("/notify-status", r.handler.GetSubscriptionNotifyStatus)
    group.Put("/:id/notify", r.handler.UpdateStreamNotifyStatus)
    group.Post("/notify/batch", r.handler.BatchUpdateNotifyStatus)
}
```

**依赖关系**: Task 6.4
**验收标准**:
- [ ] 提醒状态路由
- [ ] 批量操作路由
- [ ] 认证保护

### Task 6.9: 创建路由模块
**文件路径**: `internal/interface/http/route/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
导出路由 FX 模块：

```go
package route

import "go.uber.org/fx"

var Module = fx.Module(
    "route",
    // 路由依赖
)
```

**依赖关系**: Task 6.5-6.8
**验收标准**:
- [ ] Module 正确导出
- [ ] 可被接口层依赖

### Task 6.10: 创建接口层模块
**文件路径**: `internal/interface/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
导出接口层 FX 模块：

```go
package interface

import "go.uber.org/fx"

var Module = fx.Module(
    "interface",
    // 接口层依赖
)
```

**依赖关系**: Task 6.1-6.9
**验收标准**:
- [ ] Module 正确导出
- [ ] 可被应用层依赖

---

## 阶段 7: 应用层集成 (Priority: 3)

### Task 7.1: 更新领域模块
**文件路径**: `internal/domain/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
添加领域层 FX 依赖：

```go
package domain

import (
    "go.uber.org/fx"
    "github.com/ryuyb/fusion/internal/domain/service"
    "github.com/ryuyb/fusion/internal/domain/repository"
)

var Module = fx.Module(
    "domain",
    fx.Provide(
        // 提供服务
        service.NewNotificationService,
        service.NewNotificationPreferenceService,
        service.NewNotificationHistoryService,
        service.NewNotificationChannelService,
        service.NewSubscriptionNotificationService,
        service.NewChannelService,
        // 提供仓储接口
        repository.NewNotificationUserPreferenceRepository,
        repository.NewNotificationHistoryRepository,
        repository.NewNotificationChannelRepository,
    ),
)
```

**依赖关系**: Task 3.14
**验收标准**:
- [ ] 依赖注入配置
- [ ] 构造函数注册
- [ ] 生命周期管理

### Task 7.2: 更新应用模块
**文件路径**: `internal/application/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
添加应用层 FX 依赖：

```go
package application

import (
    "go.uber.org/fx"
    "github.com/ryuyb/fusion/internal/application/service"
)

var Module = fx.Module(
    "application",
    fx.Provide(
        service.NewNotificationService,
        service.NewNotificationPreferenceService,
        service.NewNotificationChannelService,
        service.NewSubscriptionNotificationService,
    ),
)
```

**依赖关系**: Task 4.7
**验收标准**:
- [ ] 依赖注入配置
- [ ] 服务层注册
- [ ] 生命周期管理

### Task 7.3: 更新基础设施模块
**文件路径**: `internal/infrastructure/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
添加基础设施层 FX 依赖：

```go
package infrastructure

import (
    "go.uber.org/fx"
    "github.com/ryuyb/fusion/internal/infrastructure/repository"
    "github.com/ryuyb/fusion/internal/infrastructure/sender"
)

var Module = fx.Module(
    "infrastructure",
    fx.Provide(
        repository.NewNotificationUserPreferenceRepository,
        repository.NewNotificationHistoryRepository,
        repository.NewNotificationChannelRepository,
        sender.ProvideEmailSender,
        sender.ProvidePushSender,
        sender.ProvideWebhookSender,
        sender.ProvideSlackSender,
        sender.ProvideDiscordSender,
        sender.ProvideSMSsender,
    ),
    fx.Invoke(sender.RegisterChannels),
)
```

**依赖关系**: Task 5.13
**验收标准**:
- [ ] 依赖注入配置
- [ ] 仓储实现注册
- [ ] 发送器注册

### Task 7.4: 更新接口模块
**文件路径**: `internal/interface/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
添加接口层 FX 依赖：

```go
package interface

import (
    "go.uber.org/fx"
    "github.com/ryuyb/fusion/internal/interface/http/handler"
    "github.com/ryuyb/fusion/internal/interface/http/route"
)

var Module = fx.Module(
    "interface",
    fx.Provide(
        handler.NewNotificationHandler,
        handler.NewNotificationPreferenceHandler,
        handler.NewNotificationChannelHandler,
        handler.NewSubscriptionNotificationHandler,
        route.NewNotificationRoute,
        route.NewNotificationPreferenceRoute,
        route.NewNotificationChannelRoute,
        route.NewSubscriptionNotificationRoute,
    ),
)
```

**依赖关系**: Task 6.10
**验收标准**:
- [ ] 依赖注入配置
- [ ] 处理器注册
- [ ] 路由注册

### Task 7.5: 更新应用模块
**文件路径**: `internal/app/module.go`
**关联需求**: 无
**关联设计**: 无

#### 实现细节
组合所有模块：

```go
package app

import (
    "go.uber.org/fx"
    "github.com/ryuyb/fusion/internal/domain"
    "github.com/ryuyb/fusion/internal/application"
    "github.com/ryuyb/fusion/internal/infrastructure"
    "github.com/ryuyb/fusion/internal/interface"
)

var Module = fx.Module(
    "app",
    domain.Module,
    application.Module,
    infrastructure.Module,
    interface.Module,
)
```

**依赖关系**: Task 7.1-7.4
**验收标准**:
- [ ] 所有模块组合
- [ ] 依赖注入链
- [ ] 生命周期管理

---

## 阶段 8: 测试 (Priority: 4)

### Task 8.1: 编写单元测试 - Domain 层
**文件路径**: `internal/domain/service/notification_service_test.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.8

#### 实现细节
测试通知服务：

```go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/ryuyb/fusion/internal/domain/repository"
    "go.uber.org/zap"
)

func TestNotificationService_ProcessStreamStartEvent(t *testing.T) {
    mockRepo := new(repository.MockNotificationRepository)
    mockPrefRepo := new(repository.MockNotificationUserPreferenceRepository)
    mockHistoryRepo := new(repository.MockNotificationHistoryRepository)
    mockChannelSvc := new(MockChannelService)
    logger := zap.NewNop()

    svc := NewNotificationService(
        mockRepo,
        mockPrefRepo,
        mockHistoryRepo,
        mockChannelSvc,
        nil, // eventBus
        logger,
    )

    event := StreamStartEvent{
        StreamerID:   1,
        StreamerName: "Test Streamer",
        Platform:     "test",
        Title:        "Test Stream",
        URL:          "http://test.com",
    }

    // 测试逻辑
    err := svc.ProcessStreamStartEvent(context.Background(), event)
    assert.NoError(t, err)
}
```

**依赖关系**: Task 3.8
**验收标准**:
- [ ] Mock 依赖
- [ ] 测试用例完整
- [ ] 断言正确

### Task 8.2: 编写集成测试 - Repository
**文件路径**: `test/integration/repository/notification_test.go`
**关联需求**: 需求文档 2.3
**关联设计**: 设计文档 2.8

#### 实现细节
测试仓储层：

```go
package repository

import (
    "context"
    "testing"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent"
    "github.com/ryuyb/fusion/internal/infrastructure/database/ent/notificationuserpreference"
    "entgo.io/ent/dialect/sql"
    "entgo.io/ent/dialect/sql/sqlgraph"
    "github.com/ryuyb/fusion/test/integration/testutil"
)

func TestNotificationUserPreferenceRepository_Create(t *testing.T) {
    client := testutil.EntClient(t)
    defer client.Close()

    repo := NewNotificationUserPreferenceRepository(client)
    pref := &entity.NotificationUserPreference{
        UserID: 1,
        QuietHoursEnabled: true,
    }

    err := repo.Create(context.Background(), pref)
    assert.NoError(t, err)
    assert.True(t, pref.ID > 0)
}
```

**依赖关系**: Task 5.1
**验收标准**:
- [ ] 使用 enttest
- [ ] 测试数据库
- [ ] 事务测试

### Task 8.3: 编写 API 测试 - Handler
**文件路径**: `test/integration/http/notification_test.go`
**关联需求**: 需求文档 2.4
**关联设计**: 设计文档 2.8

#### 实现细节
测试 HTTP 处理器：

```go
package http

import (
    "testing"
    "github.com/gofiber/fiber/v3"
    "github.com/stretchr/testify/assert"
    "github.com/ryuyb/fusion/test/integration/testutil"
)

func TestNotificationHandler_GetHistory(t *testing.T) {
    app := testutil.SetupApp(t)
    client := testutil.EntClient(t)

    // 设置测试数据
    token := testutil.GenerateJWT(t, 1)

    req := httptest.NewRequest("GET", "/api/notifications/history?limit=20&offset=0", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
```

**依赖关系**: Task 6.1
**验收标准**:
- [ ] Fiber 测试
- [ ] JWT 认证测试
- [ ] 响应验证

### Task 8.4: 编写端到端测试
**文件路径**: `test/e2e/notification_test.go`
**关联需求**: 需求文档 2.1
**关联设计**: 设计文档 2.8

#### 实现细节
测试完整流程：

```go
package e2e

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/ryuyb/fusion/test/e2e/testutil"
)

func TestCompleteNotificationFlow(t *testing.T) {
    client := testutil.NewClient(t)
    defer client.Close()

    // 1. 用户登录
    token := testutil.Login(t, client, "test@example.com", "password")

    // 2. 配置通知偏好
    preferences := testutil.UpdatePreferences(t, client, token, &UpdatePreferencesRequest{
        QuietHoursEnabled: true,
    })
    assert.NotNil(t, preferences)

    // 3. 订阅主播
    subscription := testutil.SubscribeStreamer(t, client, token, 1)
    assert.NotNil(t, subscription)

    // 4. 开启提醒
    testutil.UpdateNotifyStatus(t, client, token, subscription.ID, true)

    // 5. 模拟主播开播
    testutil.TriggerStreamStart(t, 1)

    // 6. 验证通知发送
    histories := testutil.GetNotificationHistory(t, client, token)
    assert.Len(t, histories, 1)
}
```

**依赖关系**: Task 6.1-6.4
**验收标准**:
- [ ] 完整用户流程
- [ ] 多服务集成
- [ ] 端到端验证

---

## 阶段 9: 文档 (Priority: 5)

### Task 9.1: 生成 Swagger 文档
**文件路径**: docs/api/notification.md
**关联需求**: 需求文档 2.1-2.4
**关联设计**: 设计文档 2.6

#### 实现细节
运行 Swagger 生成：

```bash
# 生成 Swagger 文档
make generate-swagger

# 验证文档
make format-swagger
```

**依赖关系**: Task 6.1-6.4
**验收标准**:
- [ ] 文档生成成功
- [ ] API 端点完整
- [ ] 示例正确

### Task 9.2: 编写开发文档
**文件路径**: docs/development/notification-system.md
**关联需求**: 无
**关联设计**: 无

#### 实现细节
编写开发指南：

```markdown
# 通知系统开发指南

## 概述
通知系统支持多渠道通知，包括邮件、推送、Webhook、Slack、Discord、短信等。

## 架构
- Domain: 业务逻辑
- Application: 应用服务
- Infrastructure: 外部集成
- Interface: HTTP 接口

## 使用指南
1. 配置通知渠道
2. 设置用户偏好
3. 管理订阅提醒
4. 查看通知历史

## 扩展渠道
要添加新的通知渠道，需要：
1. 实现 BaseSender 接口
2. 注册到 ChannelRegistry
3. 添加配置支持
4. 编写测试
```

**依赖关系**: Task 5.4
**验收标准**:
- [ ] 架构说明清晰
- [ ] 使用示例完整
- [ ] 扩展指南详细

### Task 9.3: 编写部署文档
**文件路径**: docs/deployment/notification.md
**关联需求**: 无
**关联设计**: 无

#### 实现细节
编写部署指南：

```markdown
# 通知系统部署指南

## 配置要求
- PostgreSQL: 用户偏好、历史记录
- Redis: 通知队列（可选）
- SMTP: 邮件发送
- VAPID: 推送通知

## 环境变量
- FUSION_NOTIFICATION_EMAIL_SMTP_HOST
- FUSION_NOTIFICATION_PUSH_VAPID_PUBLIC_KEY
- FUSION_NOTIFICATION_WEBHOOK_SIGNATURE_SECRET

## 监控指标
- 通知发送成功率
- 通知延迟
- 队列长度
- 错误率
```

**依赖关系**: Task 1.1
**验收标准**:
- [ ] 配置说明完整
- [ ] 部署步骤清晰
- [ ] 监控指标定义

---

## 阶段 10: 验收 (Priority: 5)

### Task 10.1: 功能验收测试
**文件路径**: test/acceptance/notification_test.go
**关联需求**: 需求文档 2.1-2.4
**关联设计**: 设计文档 2.8

#### 实现细节
执行验收测试：

```go
package acceptance

import (
    "testing"
    "github.com/cucumber/godog"
)

func TestNotificationFeature(t *testing.T) {
    suite := godog.TestSuite{
        ScenarioInitializer: InitializeScenario,
        Options: &godog.Options{
            Format: "pretty",
            Paths:  []string{"features"},
        },
    }

    suite.Run(t)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Before(func(ctx context.Context, sc *godog.Scenario) error {
        // 初始化测试环境
        return nil
    })

    ctx.Step(`^I have subscribed to streamer "([^"]*)"$`, iHaveSubscribedToStreamer)
    ctx.Step(`^I enabled stream notifications$`, iEnabledStreamNotifications)
    ctx.Step(`^The streamer starts streaming$`, theStreamerStartsStreaming)
    ctx.Step(`^I should receive a notification$`, iShouldReceiveANotification)
}
```

**依赖关系**: Task 8.4
**验收标准**:
- [ ] 所有用户故事验证
- [ ] 验收标准满足
- [ ] EARS 格式验证

### Task 10.2: 性能测试
**文件路径**: test/performance/notification_test.go
**关联需求**: 需求文档 4.1
**关联设计**: 设计文档 2.8

#### 实现细节
执行性能测试：

```go
package performance

import (
    "testing"
    "github.com/ory/dockertest/v3"
)

func BenchmarkNotificationSend(b *testing.B) {
    pool, err := dockertest.NewPool("")
    if err != nil {
        b.Fatalf("Could not connect to Docker: %v", err)
    }

    // 启动测试环境
    // ...

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 发送通知
        err := sendNotification()
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**依赖关系**: Task 5.5-5.7
**验收标准**:
- [ ] 响应时间 < 200ms
- [ ] 并发 1000 RPS
- [ ] 可用性 99.9%

### Task 10.3: 安全测试
**文件路径**: test/security/notification_test.go
**关联需求**: 需求文档 4.3
**关联设计**: 设计文档 2.8

#### 实现细节
执行安全测试：

```go
package security

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNotification_Security(t *testing.T) {
    t.Run("Webhook Signature Validation", func(t *testing.T) {
        // 测试签名验证
    })

    t.Run("Unauthorized Access", func(t *testing.T) {
        // 测试未授权访问
        req := httptest.NewRequest("GET", "/api/notifications/history", nil)
        resp := executeRequest(req)
        assert.Equal(t, 401, resp.StatusCode)
    })

    t.Run("Rate Limiting", func(t *testing.T) {
        // 测试频率限制
    })
}
```

**依赖关系**: Task 6.1
**验收标准**:
- [ ] 认证保护
- [ ] 签名验证
- [ ] 频率限制

---

## 实施优先级总结

### Priority 1 (关键路径)
- 阶段 1: 配置扩展
- 阶段 2: 数据模型

### Priority 2 (核心业务)
- 阶段 3: Domain 层
- 阶段 4: Application 层

### Priority 3 (外部集成)
- 阶段 5: Infrastructure 层
- 阶段 6: Interface 层

### Priority 4 (质量保证)
- 阶段 7: 应用层集成
- 阶段 8: 测试

### Priority 5 (交付)
- 阶段 9: 文档
- 阶段 10: 验收

---

## 任务依赖关系图

```
Phase 1 (Config) ─────────┐
                            │
Phase 2 (Schemas) ─────────┼─── Phase 3 (Domain) ──── Phase 4 (Application)
                            │         │                     │
Phase 5 (Infrastructure) ──┘         │                     │
                                       │                     │
Phase 6 (Interface) ──────────────────┘                     │
                                                               │
Phase 7 (Integration) ────────────────────────────────────────┘
                                                               │
Phase 8 (Tests) ──────────────────────────────────────────────┘
                                                               │
Phase 9 (Docs) ───────────────────────────────────────────────┘
                                                               │
Phase 10 (Acceptance) ─────────────────────────────────────────┘
```

---

## 验收清单

### 架构一致性
- [ ] 严格遵循 Clean Architecture
- [ ] 分层依赖正确
- [ ] 无循环依赖
- [ ] 接口隔离原则

### 功能完整性
- [ ] 所有用户故事实现
- [ ] 所有业务规则满足
- [ ] 所有非功能性需求达标

### 质量标准
- [ ] 代码覆盖率 > 80%
- [ ] API 文档完整
- [ ] 性能测试通过
- [ ] 安全测试通过

### 交付标准
- [ ] 单元测试全部通过
- [ ] 集成测试全部通过
- [ ] 端到端测试通过
- [ ] 文档完整

---

**文档版本**: v1.0
**创建日期**: 2025-11-09
**最后更新**: 2025-11-09
**维护者**: Fusion 技术团队
**状态**: 实施阶段