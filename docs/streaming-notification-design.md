# 直播平台关注与开播提醒功能设计方案

## 功能概述

设计一个可扩展的直播平台监控系统，支持用户关注多个平台的主播，并通过多种推送渠道接收开播提醒。

## 核心设计要点

### 1. Domain 实体设计 (5个核心实体)

#### 1.1 Platform (直播平台)

- **字段**：ID, 名称, 类型(douyu/huya/bilibili), 配置(API密钥等), 状态, 轮询间隔, 创建时间, 更新时间
- **特性**：支持动态添加新平台，配置使用 JSON 存储，支持不同平台的不同配置项

#### 1.2 Streamer (主播)

- **字段**：ID, 平台ID, 平台主播ID, 主播名, 头像URL, 简介, 直播间URL, 最后检查时间, 当前在线状态, 最后上线时间, 创建时间, 更新时间
- **关联**：belongs_to Platform
- **索引**：(platform_id, platform_streamer_id) 唯一

#### 1.3 UserFollowing (用户关注)

- **字段**：ID, 用户ID, 主播ID, 是否开启提醒, 最后通知时间, 创建时间, 更新时间
- **关联**：belongs_to User, belongs_to Streamer
- **索引**：(user_id, streamer_id) 唯一, streamer_id, (user_id, notification_enabled)

#### 1.4 NotificationChannel (推送渠道)

- **字段**：ID, 用户ID, 渠道类型(email/webhook/telegram/discord/feishu), 渠道名称, 配置(JSON), 是否启用, 优先级, 创建时间, 更新时间
- **关联**：belongs_to User
- **特性**：支持每用户配置多个渠道，优先级用于控制发送顺序（数字越小优先级越高）

#### 1.5 NotificationRule (提醒规则)

- **字段**：ID, 用户ID, 规则类型(silent_period/rate_limit/content_filter), 规则名称, 配置(JSON), 是否启用, 创建时间, 更新时间
- **关联**：belongs_to User
- **规则配置示例**：
  - silent_period: `{"start_hour": 23, "end_hour": 8}`
  - rate_limit: `{"interval_hours": 2}`
  - content_filter: `{"keywords": ["英雄联盟", "LOL"], "match_mode": "any"}`

### 2. 统一接口设计

#### 2.1 StreamingPlatformProvider (直播平台接口)

```go
type StreamingPlatformProvider interface {
    // 获取平台类型
    GetPlatformType() PlatformType

    // 获取主播详细信息
    FetchStreamerInfo(ctx context.Context, platformStreamerId string) (*StreamerInfo, error)

    // 检查直播状态
    CheckLiveStatus(ctx context.Context, platformStreamerId string) (*LiveStatus, error)

    // 批量检查直播状态（性能优化）
    BatchCheckLiveStatus(ctx context.Context, platformStreamerIds []string) (map[string]*LiveStatus, error)

    // 验证平台配置
    ValidateConfiguration(config map[string]interface{}) error

    // 搜索主播
    SearchStreamer(ctx context.Context, keyword string) ([]*StreamerInfo, error)
}

type StreamerInfo struct {
    PlatformStreamerId string
    Name               string
    Avatar             string
    Description        string
    RoomURL            string
}

type LiveStatus struct {
    IsLive      bool
    Title       string
    GameName    string
    StartTime   time.Time
    Viewers     int
    CoverImage  string
}
```

**实现**：
- `DouyuProvider` - 斗鱼平台
- `HuyaProvider` - 虎牙平台
- `BilibiliProvider` - 哔哩哔哩直播
- 通过 `fx.Group("streaming_providers")` 自动注册
- 使用 `resty` 发起 HTTP 请求

#### 2.2 NotificationChannelProvider (推送渠道接口)

```go
type NotificationChannelProvider interface {
    // 获取渠道类型
    GetChannelType() ChannelType

    // 发送通知
    Send(ctx context.Context, channel *NotificationChannel, notification *Notification) error

    // 验证渠道配置
    ValidateConfiguration(config map[string]interface{}) error

    // 测试渠道连通性
    TestConnection(ctx context.Context, config map[string]interface{}) error
}

type Notification struct {
    Title       string
    Content     string
    StreamerName string
    StreamerAvatar string
    RoomURL     string
    CoverImage  string
    ExtraData   map[string]interface{}
}
```

**实现**：
- `EmailChannelProvider` - 邮件通知
- `WebhookChannelProvider` - Webhook 回调
- `TelegramChannelProvider` - Telegram Bot
- `DiscordChannelProvider` - Discord Webhook
- `FeishuChannelProvider` - 飞书机器人
- 通过 `fx.Group("notification_providers")` 自动注册
- 使用 `resty` 发起 HTTP 请求

### 3. 核心服务设计 (6个应用服务)

#### 3.1 PlatformService

- `Create(dto) → Platform` - 创建平台配置（管理员）
- `Update(id, dto) → Platform` - 更新平台配置
- `GetByID(id) → Platform` - 获取平台详情
- `List() → []Platform` - 获取平台列表
- `Delete(id)` - 删除平台
- `TestConnection(id) → bool` - 测试平台API连接

#### 3.2 StreamerService

- `SyncStreamerInfo(platformId, platformStreamerId) → Streamer` - 同步主播信息
- `GetByID(id) → Streamer` - 获取主播详情
- `Search(platformType, keyword) → []Streamer` - 搜索主播
- `UpdateLiveStatus(streamerId, status)` - 更新直播状态
- `GetStreamersByPlatform(platformId) → []Streamer` - 获取平台下所有主播

#### 3.3 FollowingService

- `Follow(userId, platformType, platformStreamerId) → UserFollowing` - 关注主播
- `Unfollow(userId, followingId)` - 取消关注
- `UpdateNotificationEnabled(userId, followingId, enabled)` - 开启/关闭提醒
- `GetUserFollowings(userId, filters) → []UserFollowing` - 获取关注列表（支持按平台、提醒状态筛选）
- `GetFollowersByStreamer(streamerId) → []UserFollowing` - 获取主播的关注者
- `IsFollowing(userId, streamerId) → bool` - 检查是否已关注

#### 3.4 NotificationChannelService

- `Create(userId, dto) → NotificationChannel` - 添加推送渠道
- `Update(userId, channelId, dto) → NotificationChannel` - 更新渠道
- `Delete(userId, channelId)` - 删除渠道
- `GetByID(userId, channelId) → NotificationChannel` - 获取渠道详情
- `List(userId) → []NotificationChannel` - 获取用户所有渠道
- `TestChannel(userId, channelId) → bool` - 测试渠道连通性
- `GetEnabledChannels(userId) → []NotificationChannel` - 获取用户已启用的渠道

#### 3.5 NotificationRuleService

- `Create(userId, dto) → NotificationRule` - 创建提醒规则
- `Update(userId, ruleId, dto) → NotificationRule` - 更新规则
- `Delete(userId, ruleId)` - 删除规则
- `GetByID(userId, ruleId) → NotificationRule` - 获取规则详情
- `List(userId) → []NotificationRule` - 获取用户所有规则
- `GetEnabledRules(userId) → []NotificationRule` - 获取已启用的规则
- `ShouldNotify(userId, streamer, liveStatus) → bool` - 判断是否应该发送通知

#### 3.6 LiveCheckService (核心调度服务)

- `CheckAllStreamers()` - 检查所有主播状态（定时任务调用）
- `CheckStreamersByPlatform(platformId)` - 检查指定平台的主播
- `ProcessLiveOnEvent(streamer, liveStatus)` - 处理开播事件
- `SendNotifications(streamer, liveStatus, followers)` - 发送通知给关注者

### 4. 后台任务设计（使用 gocron）

#### 4.1 调度器配置

```go
// 使用 go-co-op/gocron/v2
scheduler := gocron.NewScheduler()

// 配置任务
job, err := scheduler.NewJob(
    gocron.DurationJob(1 * time.Minute), // 每分钟执行
    gocron.NewTask(liveCheckService.CheckAllStreamers),
)
```

#### 4.2 LiveStatusChecker (Cron Job)

- **调度频率**：每分钟（可通过配置调整）
- **任务逻辑**：
  1. 获取所有需要检查的主播列表（有用户关注的）
  2. 按平台分组
  3. 并发调用各平台的 `BatchCheckLiveStatus`
  4. 比对状态变化（离线→在线）
  5. 对于开播的主播，触发 `ProcessLiveOnEvent`

#### 4.3 NotificationDispatcher

- 查找该主播的所有开启提醒的关注者
- 对每个用户：
  1. 加载用户的提醒规则
  2. 应用规则判断是否发送（静默时段/频率限制/内容过滤）
  3. 获取用户启用的推送渠道（按优先级排序）
  4. 依次尝试发送通知（失败则尝试下一个渠道）
  5. 更新 `last_notified_at` 时间戳

### 5. 数据库设计要点

#### EntGO Schema 设计

**Platform Schema:**
```go
Fields:
- id: Int64 (unique, immutable)
- name: String (not empty)
- platform_type: Enum (douyu/huya/bilibili, immutable)
- config: JSON (平台API配置)
- status: Enum (active/inactive, default: active)
- poll_interval: Int (轮询间隔秒数, default: 60)
- created_at, updated_at

Indexes:
- platform_type (unique)
```

**Streamer Schema:**
```go
Fields:
- id: Int64 (unique, immutable)
- platform_id: Int64 (not null)
- platform_streamer_id: String (not empty)
- name: String (not empty)
- avatar: String
- description: String
- room_url: String
- last_checked_at: Time
- is_live: Bool (default: false)
- last_live_at: Time
- created_at, updated_at

Edges:
- platform: Many-to-one with Platform

Indexes:
- (platform_id, platform_streamer_id) unique
- is_live
- last_checked_at
```

**UserFollowing Schema:**
```go
Fields:
- id: Int64 (unique, immutable)
- user_id: Int64 (not null)
- streamer_id: Int64 (not null)
- notification_enabled: Bool (default: true)
- last_notified_at: Time
- created_at, updated_at

Edges:
- user: Many-to-one with User
- streamer: Many-to-one with Streamer

Indexes:
- (user_id, streamer_id) unique
- streamer_id
- (user_id, notification_enabled)
```

**NotificationChannel Schema:**
```go
Fields:
- id: Int64 (unique, immutable)
- user_id: Int64 (not null)
- channel_type: Enum (email/webhook/telegram/discord/feishu)
- name: String (not empty)
- config: JSON (渠道配置)
- is_enabled: Bool (default: true)
- priority: Int (default: 0)
- created_at, updated_at

Edges:
- user: Many-to-one with User

Indexes:
- (user_id, is_enabled)
- (user_id, priority)
```

**NotificationRule Schema:**
```go
Fields:
- id: Int64 (unique, immutable)
- user_id: Int64 (not null)
- rule_type: Enum (silent_period/rate_limit/content_filter)
- name: String (not empty)
- config: JSON (规则配置)
- is_enabled: Bool (default: true)
- created_at, updated_at

Edges:
- user: Many-to-one with User

Indexes:
- (user_id, is_enabled)
- rule_type
```

#### 软删除支持

所有实体使用 `SoftDeleteMixin`

### 6. API 端点设计 (RESTful)

#### 6.1 关注管理

- `POST /api/v1/following` - 关注主播
  - Body: `{platform_type, platform_streamer_id}`
- `DELETE /api/v1/following/:id` - 取消关注
- `GET /api/v1/following` - 获取关注列表
  - Query: `platform_type`, `notification_enabled`, `page`, `page_size`
- `PUT /api/v1/following/:id/notification` - 开启/关闭提醒
  - Body: `{enabled: bool}`
- `GET /api/v1/following/:id` - 获取关注详情

#### 6.2 推送渠道管理

- `POST /api/v1/notification-channels` - 添加推送渠道
  - Body: `{channel_type, name, config, priority}`
- `GET /api/v1/notification-channels` - 获取渠道列表
- `GET /api/v1/notification-channels/:id` - 获取渠道详情
- `PUT /api/v1/notification-channels/:id` - 更新渠道
- `DELETE /api/v1/notification-channels/:id` - 删除渠道
- `POST /api/v1/notification-channels/:id/test` - 测试渠道
- `PUT /api/v1/notification-channels/:id/toggle` - 启用/禁用渠道

#### 6.3 提醒规则管理

- `POST /api/v1/notification-rules` - 创建规则
  - Body: `{rule_type, name, config}`
- `GET /api/v1/notification-rules` - 获取规则列表
- `GET /api/v1/notification-rules/:id` - 获取规则详情
- `PUT /api/v1/notification-rules/:id` - 更新规则
- `DELETE /api/v1/notification-rules/:id` - 删除规则
- `PUT /api/v1/notification-rules/:id/toggle` - 启用/禁用规则

#### 6.4 主播搜索与查询

- `GET /api/v1/streamers/search` - 搜索主播
  - Query: `platform_type`, `keyword`
- `GET /api/v1/streamers/:id` - 获取主播详情
- `GET /api/v1/streamers/:id/live-status` - 获取主播当前直播状态

#### 6.5 平台管理（管理员端点）

- `POST /api/v1/admin/platforms` - 添加平台
- `GET /api/v1/admin/platforms` - 获取平台列表
- `GET /api/v1/admin/platforms/:id` - 获取平台详情
- `PUT /api/v1/admin/platforms/:id` - 更新平台配置
- `DELETE /api/v1/admin/platforms/:id` - 删除平台
- `POST /api/v1/admin/platforms/:id/test` - 测试平台连接

### 7. 技术实现细节

#### 7.1 Provider 注册机制

```go
// infrastructure/streaming/module.go
var StreamingModule = fx.Module("streaming",
    fx.Provide(
        // 注册 resty 客户端
        NewRestyClient,

        // 注册各平台 Provider
        fx.Annotate(
            NewDouyuProvider,
            fx.As(new(service.StreamingPlatformProvider)),
            fx.ResultTags(`group:"streaming_providers"`),
        ),
        fx.Annotate(
            NewHuyaProvider,
            fx.As(new(service.StreamingPlatformProvider)),
            fx.ResultTags(`group:"streaming_providers"`),
        ),
        fx.Annotate(
            NewBilibiliProvider,
            fx.As(new(service.StreamingPlatformProvider)),
            fx.ResultTags(`group:"streaming_providers"`),
        ),

        // Provider Manager
        NewStreamingProviderManager,
    ),
)

// StreamingProviderManager 管理所有 Provider
type StreamingProviderManager struct {
    providers map[entity.PlatformType]service.StreamingPlatformProvider
}

func NewStreamingProviderManager(
    providers []service.StreamingPlatformProvider, // fx.Group 注入
) *StreamingProviderManager {
    pm := &StreamingProviderManager{
        providers: make(map[entity.PlatformType]service.StreamingPlatformProvider),
    }
    for _, p := range providers {
        pm.providers[p.GetPlatformType()] = p
    }
    return pm
}
```

#### 7.2 Resty 客户端配置

```go
func NewRestyClient() *resty.Client {
    client := resty.New()
    client.SetTimeout(10 * time.Second)
    client.SetRetryCount(3)
    client.SetRetryWaitTime(1 * time.Second)
    client.SetRetryMaxWaitTime(5 * time.Second)
    client.SetHeader("User-Agent", "Fusion-Streaming-Platform/1.0")

    // 添加日志中间件
    client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
        // 记录请求日志
        return nil
    })

    return client
}
```

#### 7.3 Gocron 调度器配置

```go
// infrastructure/scheduler/scheduler.go
func NewScheduler(
    lc fx.Lifecycle,
    liveCheckService *service.LiveCheckService,
    logger *zap.Logger,
) (gocron.Scheduler, error) {
    s, err := gocron.NewScheduler()
    if err != nil {
        return nil, err
    }

    lc.OnStart(func(ctx context.Context) error {
        // 定义检查任务
        _, err := s.NewJob(
            gocron.DurationJob(1*time.Minute),
            gocron.NewTask(func() {
                if err := liveCheckService.CheckAllStreamers(ctx); err != nil {
                    logger.Error("Failed to check streamers", zap.Error(err))
                }
            }),
        )
        if err != nil {
            return err
        }

        // 启动调度器
        s.Start()
        logger.Info("Scheduler started")
        return nil
    })

    lc.OnStop(func(ctx context.Context) error {
        // 停止调度器
        if err := s.Shutdown(); err != nil {
            logger.Error("Failed to shutdown scheduler", zap.Error(err))
            return err
        }
        logger.Info("Scheduler stopped")
        return nil
    })

    return s, nil
}
```

#### 7.4 配置管理

在 `config.yaml` 添加新配置段：

```yaml
scheduler:
  live_check_interval: 60s  # 直播状态检查间隔

streaming:
  default_poll_interval: 60  # 默认轮询间隔（秒）
  batch_size: 50  # 批量检查主播数量

notification:
  concurrent_send: 10  # 并发发送通知数
  retry_count: 3  # 重试次数
```

#### 7.5 通知去重与频率控制

```go
// NotificationRuleService.ShouldNotify
func (s *NotificationRuleService) ShouldNotify(
    ctx context.Context,
    userId int64,
    following *entity.UserFollowing,
    liveStatus *LiveStatus,
) (bool, error) {
    rules, err := s.ruleRepo.FindEnabledByUserId(ctx, userId)
    if err != nil {
        return false, err
    }

    for _, rule := range rules {
        switch rule.RuleType {
        case entity.RuleTypeSilentPeriod:
            // 检查是否在静默时段
            if isInSilentPeriod(rule.Config) {
                return false, nil
            }

        case entity.RuleTypeRateLimit:
            // 检查通知频率限制
            if !checkRateLimit(following.LastNotifiedAt, rule.Config) {
                return false, nil
            }

        case entity.RuleTypeContentFilter:
            // 检查内容过滤规则
            if !matchContentFilter(liveStatus, rule.Config) {
                return false, nil
            }
        }
    }

    return true, nil
}
```

### 8. 扩展性设计

#### 添加新直播平台（3步）：

1. 在 `internal/infrastructure/streaming/` 创建 `xxx_provider.go`
2. 实现 `StreamingPlatformProvider` 接口（使用 resty 发起请求）
3. 在 `internal/infrastructure/streaming/module.go` 注册 Provider

#### 添加新推送渠道（3步）：

1. 在 `internal/infrastructure/notification/` 创建 `xxx_provider.go`
2. 实现 `NotificationChannelProvider` 接口（使用 resty 发起请求）
3. 在 `internal/infrastructure/notification/module.go` 注册 Provider

### 9. 文件结构

```
internal/
├── domain/
│   ├── entity/
│   │   ├── platform.go                    # 平台实体
│   │   ├── streamer.go                    # 主播实体
│   │   ├── user_following.go              # 关注关系实体
│   │   ├── notification_channel.go        # 推送渠道实体
│   │   └── notification_rule.go           # 提醒规则实体
│   ├── repository/
│   │   ├── platform_repository.go         # 平台仓储接口
│   │   ├── streamer_repository.go         # 主播仓储接口
│   │   ├── following_repository.go        # 关注仓储接口
│   │   ├── channel_repository.go          # 渠道仓储接口
│   │   └── rule_repository.go             # 规则仓储接口
│   └── service/
│       ├── streaming_platform_provider.go # 直播平台Provider接口
│       └── notification_channel_provider.go # 推送渠道Provider接口
│
├── infrastructure/
│   ├── database/
│   │   └── schema/
│   │       ├── platform.go                # Platform Ent schema
│   │       ├── streamer.go                # Streamer Ent schema
│   │       ├── user_following.go          # UserFollowing Ent schema
│   │       ├── notification_channel.go    # NotificationChannel Ent schema
│   │       └── notification_rule.go       # NotificationRule Ent schema
│   ├── repository/
│   │   ├── platform_repository.go         # 平台仓储实现
│   │   ├── streamer_repository.go         # 主播仓储实现
│   │   ├── following_repository.go        # 关注仓储实现
│   │   ├── channel_repository.go          # 渠道仓储实现
│   │   └── rule_repository.go             # 规则仓储实现
│   ├── streaming/
│   │   ├── module.go                      # streaming模块注册
│   │   ├── provider_manager.go            # Provider管理器
│   │   ├── douyu_provider.go              # 斗鱼Provider实现
│   │   ├── huya_provider.go               # 虎牙Provider实现
│   │   └── bilibili_provider.go           # B站Provider实现
│   ├── notification/
│   │   ├── module.go                      # notification模块注册
│   │   ├── provider_manager.go            # Provider管理器
│   │   ├── email_provider.go              # 邮件Provider实现
│   │   ├── webhook_provider.go            # Webhook Provider实现
│   │   ├── telegram_provider.go           # Telegram Provider实现
│   │   ├── discord_provider.go            # Discord Provider实现
│   │   └── feishu_provider.go             # 飞书Provider实现
│   └── scheduler/
│       ├── module.go                      # scheduler模块注册
│       └── scheduler.go                   # gocron调度器配置
│
├── application/
│   └── service/
│       ├── platform_service.go            # 平台服务
│       ├── streamer_service.go            # 主播服务
│       ├── following_service.go           # 关注服务
│       ├── notification_channel_service.go # 推送渠道服务
│       ├── notification_rule_service.go   # 提醒规则服务
│       └── live_check_service.go          # 直播检查服务-核心调度
│
└── interface/
    └── http/
        ├── dto/
        │   ├── request/
        │   │   ├── following_request.go   # 关注相关请求DTO
        │   │   ├── channel_request.go     # 推送渠道请求DTO
        │   │   ├── rule_request.go        # 提醒规则请求DTO
        │   │   ├── platform_request.go    # 平台请求DTO
        │   │   └── streamer_request.go    # 主播请求DTO
        │   └── response/
        │       ├── following_response.go  # 关注响应DTO
        │       ├── channel_response.go    # 渠道响应DTO
        │       ├── rule_response.go       # 规则响应DTO
        │       ├── platform_response.go   # 平台响应DTO
        │       └── streamer_response.go   # 主播响应DTO
        ├── handler/
        │   ├── following_handler.go       # 关注处理器
        │   ├── notification_channel_handler.go # 推送渠道处理器
        │   ├── notification_rule_handler.go    # 提醒规则处理器
        │   ├── streamer_handler.go        # 主播处理器
        │   └── platform_handler.go        # 平台处理器-管理员
        └── route/
            ├── following_route.go         # 关注路由
            ├── notification_route.go      # 通知相关路由
            ├── streamer_route.go          # 主播路由
            └── platform_route.go          # 平台路由-管理员
```

### 10. 实现步骤（分阶段）

#### Phase 1: 基础设施与数据模型 (2-3天)

1. 创建 5 个 Domain 实体
2. 创建 5 个 Repository 接口
3. 创建 5 个 Ent schemas
4. 生成 Ent 代码
5. 实现 5 个 Repository 实现
6. 定义 2 个 Provider 接口

#### Phase 2: 直播平台集成 (3-4天)

7. 配置 Resty 客户端
8. 实现 ProviderManager
9. 实现 DouyuProvider (使用 resty)
10. 实现 HuyaProvider (使用 resty)
11. 实现 BilibiliProvider (使用 resty)
12. 实现 PlatformService
13. 实现 StreamerService
14. 创建 DTOs (platform_request.go, platform_response.go, streamer_response.go)
15. 创建平台管理 API (handler + route)
16. 创建主播搜索 API

#### Phase 3: 关注功能 (2天)

17. 实现 FollowingService
18. 创建 DTOs (following_request.go, following_response.go)
19. 创建关注管理 API (handler + route)
20. 实现关注列表查询与筛选

#### Phase 4: 推送渠道集成 (3-4天)

21. 实现 NotificationProviderManager
22. 实现 EmailChannelProvider
23. 实现 WebhookChannelProvider (使用 resty)
24. 实现 TelegramChannelProvider (使用 resty)
25. 实现 DiscordChannelProvider (使用 resty)
26. 实现 FeishuChannelProvider (使用 resty)
27. 实现 NotificationChannelService
28. 创建 DTOs (channel_request.go, channel_response.go)
29. 创建推送渠道管理 API (handler + route)

#### Phase 5: 提醒规则 (2天)

30. 实现 NotificationRuleService
31. 实现规则判断逻辑（静默时段/频率限制/内容过滤）
32. 创建 DTOs (rule_request.go, rule_response.go)
33. 创建提醒规则管理 API (handler + route)

#### Phase 6: 核心调度系统 (3-4天)

34. 配置 gocron 调度器
35. 实现 LiveCheckService 核心逻辑
36. 实现批量检查主播状态
37. 实现开播事件处理
38. 实现通知分发逻辑
39. 集成调度器到应用生命周期

#### Phase 7: 测试与优化 (2-3天)

40. 添加单元测试
41. 并发优化（使用 errgroup）
42. 监控与日志完善
43. 错误处理优化
44. API 文档（Swagger）

**总计：约 17-20 天开发周期**

### 11. 关键技术选型与依赖

```go
// go.mod 新增依赖
require (
    github.com/go-co-op/gocron/v2 v2.x.x        // 任务调度
    github.com/go-resty/resty/v2 v2.x.x         // HTTP 客户端
    gopkg.in/gomail.v2 v2.0.0                    // 邮件发送
    github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.x.x  // Telegram
    golang.org/x/sync v0.x.x                     // errgroup 并发控制
)
```

### 12. 核心流程图

#### 开播检测与通知流程：

```
[Gocron定时任务 (每分钟)]
    ↓
[LiveCheckService.CheckAllStreamers()]
    ↓
[查询所有有关注者的主播] → [按平台分组]
    ↓
[并发调用各平台Provider.BatchCheckLiveStatus()]
    ↓
[比对状态变化] → [发现开播事件]
    ↓
[ProcessLiveOnEvent(streamer, liveStatus)]
    ↓
[查询该主播的所有关注者 (notification_enabled=true)]
    ↓
[对每个用户并发处理]
    ↓
[NotificationRuleService.ShouldNotify()] → 应用规则
    │
    ├─ [静默时段检查]
    ├─ [频率限制检查]
    └─ [内容过滤检查]
    ↓
[获取用户启用的推送渠道 (按优先级排序)]
    ↓
[依次尝试发送通知 (失败则尝试下一个)]
    ↓
[更新 last_notified_at]
```

## 总结

本方案完全遵循项目现有的 Clean Architecture 架构，通过 Provider 模式实现高度可扩展性：

### 关键特性：

- ✅ **DTOs 位置优化**：所有 request 和 response DTOs 统一放在 `internal/interface/http/dto/` 下，与现有项目结构保持一致
- ✅ **使用 gocron**：轻量级的 Go 任务调度库，支持灵活的任务配置
- ✅ **使用 resty**：功能强大的 HTTP 客户端，支持重试、中间件、链式调用
- ✅ **Provider 模式**：直播平台和推送渠道均可通过实现接口快速扩展
- ✅ **Clean Architecture**：严格遵循分层架构，依赖倒置，便于测试和维护
- ✅ **Dependency Injection**：使用 uber-go/fx 实现依赖注入，自动装配组件
- ✅ **软删除支持**：所有实体支持软删除，保留历史数据
- ✅ **并发优化**：批量检查、并发通知，提高系统性能
- ✅ **智能提醒**：支持静默时段、频率限制、内容过滤等高级规则