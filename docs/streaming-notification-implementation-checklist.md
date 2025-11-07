# 直播平台关注与开播提醒功能 - 实现清单

## 项目概览

**预计开发周期**: 17-20 天
**开始日期**: ___________
**预计完成日期**: ___________

---

## Phase 1: 基础设施与数据模型 (2-3天)

### 1.1 Domain 实体层

- [x] **创建 Platform 实体** (`internal/domain/entity/platform.go`)
  - [x] 定义 Platform 结构体（ID, Name, PlatformType, Config, Status, PollInterval, CreatedAt, UpdatedAt）
  - [x] 定义 PlatformType 枚举（Douyu, Huya, Bilibili）
  - [x] 定义 PlatformStatus 枚举（Active, Inactive）
  - [x] 添加工厂方法 `CreatePlatform()`
  - [x] 添加更新方法 `Update()`

- [x] **创建 Streamer 实体** (`internal/domain/entity/streamer.go`)
  - [x] 定义 Streamer 结构体（ID, PlatformID, PlatformStreamerID, Name, Avatar, Description, RoomURL, LastCheckedAt, IsLive, LastLiveAt, CreatedAt, UpdatedAt）
  - [x] 添加工厂方法 `CreateStreamer()`
  - [x] 添加更新方法 `Update()`
  - [x] 添加直播状态更新方法 `UpdateLiveStatus()`

- [x] **创建 UserFollowing 实体** (`internal/domain/entity/user_following.go`)
  - [x] 定义 UserFollowing 结构体（ID, UserID, StreamerID, NotificationEnabled, LastNotifiedAt, CreatedAt, UpdatedAt）
  - [x] 添加工厂方法 `CreateFollowing()`
  - [x] 添加提醒开关方法 `ToggleNotification()`
  - [x] 添加更新最后通知时间方法 `UpdateLastNotifiedAt()`

- [x] **创建 NotificationChannel 实体** (`internal/domain/entity/notification_channel.go`)
  - [x] 定义 NotificationChannel 结构体（ID, UserID, ChannelType, Name, Config, IsEnabled, Priority, CreatedAt, UpdatedAt）
  - [x] 定义 ChannelType 枚举（Email, Webhook, Telegram, Discord, Feishu）
  - [x] 添加工厂方法 `CreateChannel()`
  - [x] 添加更新方法 `Update()`
  - [x] 添加启用/禁用方法 `Toggle()`

- [x] **创建 NotificationRule 实体** (`internal/domain/entity/notification_rule.go`)
  - [x] 定义 NotificationRule 结构体（ID, UserID, RuleType, Name, Config, IsEnabled, CreatedAt, UpdatedAt）
  - [x] 定义 RuleType 枚举（SilentPeriod, RateLimit, ContentFilter）
  - [x] 添加工厂方法 `CreateRule()`
  - [x] 添加更新方法 `Update()`
  - [x] 添加启用/禁用方法 `Toggle()`

### 1.2 Repository 接口层

- [x] **创建 PlatformRepository 接口** (`internal/domain/repository/platform_repository.go`)
  - [x] `Create(ctx, platform) → platform, error`
  - [x] `Update(ctx, platform) → platform, error`
  - [x] `FindByID(ctx, id) → platform, error`
  - [x] `FindByType(ctx, platformType) → platform, error`
  - [x] `List(ctx) → []platform, error`
  - [x] `Delete(ctx, id) → error`

- [x] **创建 StreamerRepository 接口** (`internal/domain/repository/streamer_repository.go`)
  - [x] `Create(ctx, streamer) → streamer, error`
  - [x] `Update(ctx, streamer) → streamer, error`
  - [x] `FindByID(ctx, id) → streamer, error`
  - [x] `FindByPlatformAndStreamerID(ctx, platformID, platformStreamerID) → streamer, error`
  - [x] `FindByPlatform(ctx, platformID) → []streamer, error`
  - [x] `FindAllWithFollowers(ctx) → []streamer, error` // 获取有关注者的主播
  - [x] `UpdateLiveStatus(ctx, streamerID, isLive, lastLiveAt) → error`

- [x] **创建 FollowingRepository 接口** (`internal/domain/repository/following_repository.go`)
  - [x] `Create(ctx, following) → following, error`
  - [x] `Update(ctx, following) → following, error`
  - [x] `FindByID(ctx, id) → following, error`
  - [x] `FindByUserAndStreamer(ctx, userID, streamerID) → following, error`
  - [x] `FindByUser(ctx, userID, filters) → []following, error`
  - [x] `FindByStreamer(ctx, streamerID, notificationEnabled) → []following, error`
  - [x] `Delete(ctx, id) → error`
  - [x] `UpdateLastNotifiedAt(ctx, id, time) → error`

- [x] **创建 ChannelRepository 接口** (`internal/domain/repository/channel_repository.go`)
  - [x] `Create(ctx, channel) → channel, error`
  - [x] `Update(ctx, channel) → channel, error`
  - [x] `FindByID(ctx, id) → channel, error`
  - [x] `FindByUser(ctx, userID) → []channel, error`
  - [x] `FindEnabledByUser(ctx, userID) → []channel, error` // 按优先级排序
  - [x] `Delete(ctx, id) → error`

- [x] **创建 RuleRepository 接口** (`internal/domain/repository/rule_repository.go`)
  - [x] `Create(ctx, rule) → rule, error`
  - [x] `Update(ctx, rule) → rule, error`
  - [x] `FindByID(ctx, id) → rule, error`
  - [x] `FindByUser(ctx, userID) → []rule, error`
  - [x] `FindEnabledByUser(ctx, userID) → []rule, error`
  - [x] `Delete(ctx, id) → error`

### 1.3 Provider 接口定义

- [x] **创建 StreamingPlatformProvider 接口** (`internal/domain/service/streaming_platform_provider.go`)
  - [x] 定义 `GetPlatformType()` 方法
  - [x] 定义 `FetchStreamerInfo()` 方法
  - [x] 定义 `CheckLiveStatus()` 方法
  - [x] 定义 `BatchCheckLiveStatus()` 方法
  - [x] 定义 `ValidateConfiguration()` 方法
  - [x] 定义 `SearchStreamer()` 方法
  - [x] 定义 StreamerInfo 结构体
  - [x] 定义 LiveStatus 结构体

- [x] **创建 NotificationChannelProvider 接口** (`internal/domain/service/notification_channel_provider.go`)
  - [x] 定义 `GetChannelType()` 方法
  - [x] 定义 `Send()` 方法
  - [x] 定义 `ValidateConfiguration()` 方法
  - [x] 定义 `TestConnection()` 方法
  - [x] 定义 Notification 结构体

### 1.4 Ent Schema 定义

- [x] **创建 Platform Schema** (`internal/infrastructure/database/schema/platform.go`)
  - [x] 定义字段（id, name, platform_type, config, status, poll_interval, created_at, updated_at）
  - [x] 添加 SoftDeleteMixin
  - [x] 定义索引（platform_type unique）
  - [x] 定义 Edges（has many Streamers）

- [x] **创建 Streamer Schema** (`internal/infrastructure/database/schema/streamer.go`)
  - [x] 定义字段（id, platform_id, platform_streamer_id, name, avatar, description, room_url, last_checked_at, is_live, last_live_at, created_at, updated_at）
  - [x] 添加 SoftDeleteMixin
  - [x] 定义索引（(platform_id, platform_streamer_id) unique, is_live, last_checked_at）
  - [x] 定义 Edges（belongs to Platform, has many UserFollowings）

- [x] **创建 UserFollowing Schema** (`internal/infrastructure/database/schema/user_following.go`)
  - [x] 定义字段（id, user_id, streamer_id, notification_enabled, last_notified_at, created_at, updated_at）
  - [x] 添加 SoftDeleteMixin
  - [x] 定义索引（(user_id, streamer_id) unique, streamer_id, (user_id, notification_enabled)）
  - [x] 定义 Edges（belongs to User, belongs to Streamer）

- [x] **创建 NotificationChannel Schema** (`internal/infrastructure/database/schema/notification_channel.go`)
  - [x] 定义字段（id, user_id, channel_type, name, config, is_enabled, priority, created_at, updated_at）
  - [x] 添加 SoftDeleteMixin
  - [x] 定义索引（(user_id, is_enabled), (user_id, priority)）
  - [x] 定义 Edges（belongs to User）

- [x] **创建 NotificationRule Schema** (`internal/infrastructure/database/schema/notification_rule.go`)
  - [x] 定义字段（id, user_id, rule_type, name, config, is_enabled, created_at, updated_at）
  - [x] 添加 SoftDeleteMixin
  - [x] 定义索引（(user_id, is_enabled), rule_type）
  - [x] 定义 Edges（belongs to User）

- [x] **生成 Ent 代码**
  - [x] 运行 `make generate-ent`
  - [x] 确认生成的代码无错误

### 1.5 Repository 实现层

- [ ] **实现 PlatformRepository** (`internal/infrastructure/repository/platform_repository.go`)
  - [ ] 实现所有接口方法
  - [ ] 处理数据库错误转换
  - [ ] 添加日志记录

- [ ] **实现 StreamerRepository** (`internal/infrastructure/repository/streamer_repository.go`)
  - [ ] 实现所有接口方法
  - [ ] 实现预加载 Platform 关联
  - [ ] 处理数据库错误转换
  - [ ] 添加日志记录

- [ ] **实现 FollowingRepository** (`internal/infrastructure/repository/following_repository.go`)
  - [ ] 实现所有接口方法
  - [ ] 实现预加载 User 和 Streamer 关联
  - [ ] 实现筛选逻辑（按平台、提醒状态）
  - [ ] 处理数据库错误转换
  - [ ] 添加日志记录

- [ ] **实现 ChannelRepository** (`internal/infrastructure/repository/channel_repository.go`)
  - [ ] 实现所有接口方法
  - [ ] 实现按优先级排序
  - [ ] 处理数据库错误转换
  - [ ] 添加日志记录

- [ ] **实现 RuleRepository** (`internal/infrastructure/repository/rule_repository.go`)
  - [ ] 实现所有接口方法
  - [ ] 处理数据库错误转换
  - [ ] 添加日志记录

- [ ] **创建 Repository Module** (`internal/infrastructure/repository/module.go`)
  - [ ] 添加所有 Repository 的 Provider
  - [ ] 使用 fx.Annotate 进行接口绑定

---

## Phase 2: 直播平台集成 (3-4天)

### 2.1 HTTP 客户端配置

- [x] **配置 Resty 客户端** (`internal/infrastructure/streaming/resty_client.go`)
  - [x] 创建 `NewRestyClient()` 函数
  - [x] 配置超时时间（10秒）
  - [x] 配置重试策略（3次重试，1-5秒间隔）
  - [x] 添加 User-Agent header
  - [x] 添加请求日志中间件
  - [x] 添加响应日志中间件

### 2.2 直播平台 Provider 实现

- [ ] **实现 DouyuProvider** (`internal/infrastructure/streaming/douyu_provider.go`)
  - [ ] 实现 `GetPlatformType()` 返回 Douyu
  - [ ] 实现 `FetchStreamerInfo()` - 调用斗鱼 API
  - [ ] 实现 `CheckLiveStatus()` - 检查单个主播状态
  - [ ] 实现 `BatchCheckLiveStatus()` - 批量检查（优化性能）
  - [ ] 实现 `ValidateConfiguration()` - 验证 API 配置
  - [ ] 实现 `SearchStreamer()` - 搜索主播
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [ ] **实现 HuyaProvider** (`internal/infrastructure/streaming/huya_provider.go`)
  - [ ] 实现 `GetPlatformType()` 返回 Huya
  - [ ] 实现 `FetchStreamerInfo()` - 调用虎牙 API
  - [ ] 实现 `CheckLiveStatus()` - 检查单个主播状态
  - [ ] 实现 `BatchCheckLiveStatus()` - 批量检查（优化性能）
  - [ ] 实现 `ValidateConfiguration()` - 验证 API 配置
  - [ ] 实现 `SearchStreamer()` - 搜索主播
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [x] **实现 BilibiliProvider** (`internal/infrastructure/streaming/bilibili_provider.go`)
  - [x] 实现 `GetPlatformType()` 返回 Bilibili
  - [x] 实现 `FetchStreamerInfo()` - 调用B站 API
  - [x] 实现 `CheckLiveStatus()` - 检查单个主播状态
  - [x] 实现 `BatchCheckLiveStatus()` - 批量检查（优化性能）
  - [x] 实现 `ValidateConfiguration()` - 验证 API 配置
  - [x] 实现 `SearchStreamer()` - 搜索主播
  - [x] 添加错误处理和日志
  - [x] 添加单元测试

### 2.3 Provider 管理器

- [x] **实现 StreamingProviderManager** (`internal/infrastructure/streaming/provider_manager.go`)
  - [x] 创建 ProviderManager 结构体（使用 map 存储 Provider）
  - [x] 实现 `NewStreamingProviderManager()` - 接收 fx.Group 注入的 Providers
  - [x] 实现 `GetProvider(platformType)` - 获取指定平台的 Provider
  - [x] 实现 `GetAllProviders()` - 获取所有 Providers
  - [x] 添加日志记录

- [x] **创建 Streaming Module** (`internal/infrastructure/streaming/module.go`)
  - [x] 注册 RestyClient Provider
  - [ ] 注册 DouyuProvider（使用 fx.Annotate + fx.ResultTags）
  - [ ] 注册 HuyaProvider（使用 fx.Annotate + fx.ResultTags）
  - [x] 注册 BilibiliProvider（使用 fx.Annotate + fx.ResultTags）
  - [x] 注册 StreamingProviderManager

### 2.4 应用服务实现

- [ ] **实现 PlatformService** (`internal/application/service/platform_service.go`)
  - [ ] 注入依赖（PlatformRepository, StreamingProviderManager, Logger）
  - [ ] 实现 `Create()` - 创建平台配置
  - [ ] 实现 `Update()` - 更新平台配置
  - [ ] 实现 `GetByID()` - 获取平台详情
  - [ ] 实现 `List()` - 获取平台列表
  - [ ] 实现 `Delete()` - 删除平台
  - [ ] 实现 `TestConnection()` - 测试平台 API 连接
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [ ] **实现 StreamerService** (`internal/application/service/streamer_service.go`)
  - [ ] 注入依赖（StreamerRepository, PlatformRepository, StreamingProviderManager, Logger）
  - [ ] 实现 `SyncStreamerInfo()` - 同步主播信息
  - [ ] 实现 `GetByID()` - 获取主播详情
  - [ ] 实现 `Search()` - 搜索主播（调用 Provider）
  - [ ] 实现 `UpdateLiveStatus()` - 更新直播状态
  - [ ] 实现 `GetStreamersByPlatform()` - 获取平台下所有主播
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

### 2.5 DTOs 定义

- [ ] **创建 Platform Request DTOs** (`internal/interface/http/dto/request/platform_request.go`)
  - [ ] 定义 `CreatePlatformRequest`（name, platform_type, config, poll_interval）
  - [ ] 定义 `UpdatePlatformRequest`（name, config, status, poll_interval）
  - [ ] 添加验证规则（required, enum, min, max）

- [ ] **创建 Platform Response DTOs** (`internal/interface/http/dto/response/platform_response.go`)
  - [ ] 定义 `PlatformResponse`（ID, Name, PlatformType, Config, Status, PollInterval, CreatedAt, UpdatedAt）
  - [ ] 实现转换方法 `ToPlatformResponse(entity)`

- [ ] **创建 Streamer Request DTOs** (`internal/interface/http/dto/request/streamer_request.go`)
  - [ ] 定义 `SearchStreamerRequest`（platform_type, keyword）
  - [ ] 添加验证规则

- [ ] **创建 Streamer Response DTOs** (`internal/interface/http/dto/response/streamer_response.go`)
  - [ ] 定义 `StreamerResponse`（ID, PlatformID, PlatformStreamerID, Name, Avatar, Description, RoomURL, IsLive, LastLiveAt, CreatedAt, UpdatedAt）
  - [ ] 定义 `LiveStatusResponse`（IsLive, Title, GameName, StartTime, Viewers, CoverImage）
  - [ ] 实现转换方法

### 2.6 HTTP Handler 实现

- [ ] **实现 PlatformHandler** (`internal/interface/http/handler/platform_handler.go`)
  - [ ] 注入依赖（PlatformService, Validator, Logger）
  - [ ] 实现 `Create()` - POST /api/v1/admin/platforms
  - [ ] 实现 `Update()` - PUT /api/v1/admin/platforms/:id
  - [ ] 实现 `GetByID()` - GET /api/v1/admin/platforms/:id
  - [ ] 实现 `List()` - GET /api/v1/admin/platforms
  - [ ] 实现 `Delete()` - DELETE /api/v1/admin/platforms/:id
  - [ ] 实现 `TestConnection()` - POST /api/v1/admin/platforms/:id/test
  - [ ] 添加 Swagger 注解
  - [ ] 添加错误处理

- [ ] **实现 StreamerHandler** (`internal/interface/http/handler/streamer_handler.go`)
  - [ ] 注入依赖（StreamerService, Validator, Logger）
  - [ ] 实现 `Search()` - GET /api/v1/streamers/search
  - [ ] 实现 `GetByID()` - GET /api/v1/streamers/:id
  - [ ] 实现 `GetLiveStatus()` - GET /api/v1/streamers/:id/live-status
  - [ ] 添加 Swagger 注解
  - [ ] 添加错误处理

### 2.7 路由注册

- [ ] **创建 PlatformRoute** (`internal/interface/http/route/platform_route.go`)
  - [ ] 实现 `RegisterRouters()` 方法
  - [ ] 注册所有平台管理路由（需要管理员权限）
  - [ ] 添加认证中间件

- [ ] **创建 StreamerRoute** (`internal/interface/http/route/streamer_route.go`)
  - [ ] 实现 `RegisterRouters()` 方法
  - [ ] 注册所有主播相关路由
  - [ ] 添加认证中间件（可选认证）

---

## Phase 3: 关注功能 (2天)

### 3.1 应用服务实现

- [ ] **实现 FollowingService** (`internal/application/service/following_service.go`)
  - [ ] 注入依赖（FollowingRepository, StreamerRepository, StreamerService, Logger）
  - [ ] 实现 `Follow()` - 关注主播（检查是否已关注，自动同步主播信息）
  - [ ] 实现 `Unfollow()` - 取消关注（软删除）
  - [ ] 实现 `UpdateNotificationEnabled()` - 开启/关闭提醒
  - [ ] 实现 `GetUserFollowings()` - 获取关注列表（支持筛选）
  - [ ] 实现 `GetFollowersByStreamer()` - 获取主播的关注者
  - [ ] 实现 `IsFollowing()` - 检查是否已关注
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

### 3.2 DTOs 定义

- [ ] **创建 Following Request DTOs** (`internal/interface/http/dto/request/following_request.go`)
  - [ ] 定义 `FollowRequest`（platform_type, platform_streamer_id）
  - [ ] 定义 `UpdateNotificationRequest`（enabled）
  - [ ] 定义 `ListFollowingRequest`（platform_type, notification_enabled, page, page_size）
  - [ ] 添加验证规则

- [ ] **创建 Following Response DTOs** (`internal/interface/http/dto/response/following_response.go`)
  - [ ] 定义 `FollowingResponse`（ID, UserID, Streamer, NotificationEnabled, LastNotifiedAt, CreatedAt, UpdatedAt）
  - [ ] 实现转换方法 `ToFollowingResponse(entity)`
  - [ ] 实现分页响应

### 3.3 HTTP Handler 实现

- [ ] **实现 FollowingHandler** (`internal/interface/http/handler/following_handler.go`)
  - [ ] 注入依赖（FollowingService, Validator, Logger）
  - [ ] 实现 `Follow()` - POST /api/v1/following
  - [ ] 实现 `Unfollow()` - DELETE /api/v1/following/:id
  - [ ] 实现 `List()` - GET /api/v1/following
  - [ ] 实现 `GetByID()` - GET /api/v1/following/:id
  - [ ] 实现 `UpdateNotification()` - PUT /api/v1/following/:id/notification
  - [ ] 添加 Swagger 注解
  - [ ] 添加错误处理
  - [ ] 确保用户只能操作自己的关注

### 3.4 路由注册

- [ ] **创建 FollowingRoute** (`internal/interface/http/route/following_route.go`)
  - [ ] 实现 `RegisterRouters()` 方法
  - [ ] 注册所有关注相关路由
  - [ ] 添加认证中间件（必须登录）

---

## Phase 4: 推送渠道集成 (3-4天)

### 4.1 推送渠道 Provider 实现

- [ ] **实现 EmailChannelProvider** (`internal/infrastructure/notification/email_provider.go`)
  - [ ] 实现 `GetChannelType()` 返回 Email
  - [ ] 实现 `Send()` - 使用 gomail 发送邮件
  - [ ] 实现 `ValidateConfiguration()` - 验证 SMTP 配置
  - [ ] 实现 `TestConnection()` - 测试邮件发送
  - [ ] 添加 HTML 邮件模板
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [ ] **实现 WebhookChannelProvider** (`internal/infrastructure/notification/webhook_provider.go`)
  - [ ] 实现 `GetChannelType()` 返回 Webhook
  - [ ] 实现 `Send()` - 使用 resty 发送 POST 请求
  - [ ] 实现 `ValidateConfiguration()` - 验证 URL 配置
  - [ ] 实现 `TestConnection()` - 测试 Webhook 连通性
  - [ ] 支持自定义 Headers 和 Body 格式
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [ ] **实现 TelegramChannelProvider** (`internal/infrastructure/notification/telegram_provider.go`)
  - [ ] 实现 `GetChannelType()` 返回 Telegram
  - [ ] 实现 `Send()` - 使用 Telegram Bot API 发送消息
  - [ ] 实现 `ValidateConfiguration()` - 验证 Bot Token 和 Chat ID
  - [ ] 实现 `TestConnection()` - 测试 Bot 连通性
  - [ ] 支持消息格式化（Markdown）
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [ ] **实现 DiscordChannelProvider** (`internal/infrastructure/notification/discord_provider.go`)
  - [ ] 实现 `GetChannelType()` 返回 Discord
  - [ ] 实现 `Send()` - 使用 Discord Webhook API
  - [ ] 实现 `ValidateConfiguration()` - 验证 Webhook URL
  - [ ] 实现 `TestConnection()` - 测试 Webhook 连通性
  - [ ] 支持 Embed 消息格式
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

- [ ] **实现 FeishuChannelProvider** (`internal/infrastructure/notification/feishu_provider.go`)
  - [ ] 实现 `GetChannelType()` 返回 Feishu
  - [ ] 实现 `Send()` - 使用飞书机器人 Webhook
  - [ ] 实现 `ValidateConfiguration()` - 验证 Webhook URL
  - [ ] 实现 `TestConnection()` - 测试机器人连通性
  - [ ] 支持消息卡片格式
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

### 4.2 Provider 管理器

- [ ] **实现 NotificationProviderManager** (`internal/infrastructure/notification/provider_manager.go`)
  - [ ] 创建 ProviderManager 结构体
  - [ ] 实现 `NewNotificationProviderManager()` - 接收 fx.Group 注入的 Providers
  - [ ] 实现 `GetProvider(channelType)` - 获取指定类型的 Provider
  - [ ] 实现 `GetAllProviders()` - 获取所有 Providers
  - [ ] 添加日志记录

- [ ] **创建 Notification Module** (`internal/infrastructure/notification/module.go`)
  - [ ] 注册 EmailChannelProvider（使用 fx.Annotate + fx.ResultTags）
  - [ ] 注册 WebhookChannelProvider
  - [ ] 注册 TelegramChannelProvider
  - [ ] 注册 DiscordChannelProvider
  - [ ] 注册 FeishuChannelProvider
  - [ ] 注册 NotificationProviderManager

### 4.3 应用服务实现

- [ ] **实现 NotificationChannelService** (`internal/application/service/notification_channel_service.go`)
  - [ ] 注入依赖（ChannelRepository, NotificationProviderManager, Logger）
  - [ ] 实现 `Create()` - 创建推送渠道（验证配置）
  - [ ] 实现 `Update()` - 更新渠道
  - [ ] 实现 `Delete()` - 删除渠道
  - [ ] 实现 `GetByID()` - 获取渠道详情
  - [ ] 实现 `List()` - 获取用户所有渠道
  - [ ] 实现 `TestChannel()` - 测试渠道连通性
  - [ ] 实现 `GetEnabledChannels()` - 获取已启用的渠道（按优先级）
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

### 4.4 DTOs 定义

- [ ] **创建 Channel Request DTOs** (`internal/interface/http/dto/request/channel_request.go`)
  - [ ] 定义 `CreateChannelRequest`（channel_type, name, config, priority）
  - [ ] 定义 `UpdateChannelRequest`（name, config, priority, is_enabled）
  - [ ] 添加验证规则

- [ ] **创建 Channel Response DTOs** (`internal/interface/http/dto/response/channel_response.go`)
  - [ ] 定义 `ChannelResponse`（ID, UserID, ChannelType, Name, Config, IsEnabled, Priority, CreatedAt, UpdatedAt）
  - [ ] 实现转换方法
  - [ ] 脱敏处理（隐藏敏感配置信息）

### 4.5 HTTP Handler 实现

- [ ] **实现 NotificationChannelHandler** (`internal/interface/http/handler/notification_channel_handler.go`)
  - [ ] 注入依赖（NotificationChannelService, Validator, Logger）
  - [ ] 实现 `Create()` - POST /api/v1/notification-channels
  - [ ] 实现 `Update()` - PUT /api/v1/notification-channels/:id
  - [ ] 实现 `Delete()` - DELETE /api/v1/notification-channels/:id
  - [ ] 实现 `GetByID()` - GET /api/v1/notification-channels/:id
  - [ ] 实现 `List()` - GET /api/v1/notification-channels
  - [ ] 实现 `Test()` - POST /api/v1/notification-channels/:id/test
  - [ ] 实现 `Toggle()` - PUT /api/v1/notification-channels/:id/toggle
  - [ ] 添加 Swagger 注解
  - [ ] 添加错误处理
  - [ ] 确保用户只能操作自己的渠道

### 4.6 路由注册

- [ ] **创建 NotificationRoute** (`internal/interface/http/route/notification_route.go`)
  - [ ] 实现 `RegisterRouters()` 方法
  - [ ] 注册推送渠道相关路由
  - [ ] 注册提醒规则相关路由（Phase 5 完成后）
  - [ ] 添加认证中间件（必须登录）

---

## Phase 5: 提醒规则 (2天)

### 5.1 规则判断逻辑实现

- [ ] **实现规则辅助函数** (`internal/application/service/notification_rule_helper.go`)
  - [ ] 实现 `isInSilentPeriod(config, currentTime)` - 检查是否在静默时段
  - [ ] 实现 `checkRateLimit(lastNotifiedAt, config)` - 检查频率限制
  - [ ] 实现 `matchContentFilter(liveStatus, config)` - 检查内容过滤
  - [ ] 添加单元测试（各种边界情况）

### 5.2 应用服务实现

- [ ] **实现 NotificationRuleService** (`internal/application/service/notification_rule_service.go`)
  - [ ] 注入依赖（RuleRepository, Logger）
  - [ ] 实现 `Create()` - 创建提醒规则
  - [ ] 实现 `Update()` - 更新规则
  - [ ] 实现 `Delete()` - 删除规则
  - [ ] 实现 `GetByID()` - 获取规则详情
  - [ ] 实现 `List()` - 获取用户所有规则
  - [ ] 实现 `GetEnabledRules()` - 获取已启用的规则
  - [ ] 实现 `ShouldNotify()` - 综合判断是否应该发送通知
  - [ ] 添加错误处理和日志
  - [ ] 添加单元测试

### 5.3 DTOs 定义

- [ ] **创建 Rule Request DTOs** (`internal/interface/http/dto/request/rule_request.go`)
  - [ ] 定义 `CreateRuleRequest`（rule_type, name, config）
  - [ ] 定义 `UpdateRuleRequest`（name, config, is_enabled）
  - [ ] 定义配置示例文档（各种规则类型的配置格式）
  - [ ] 添加验证规则

- [ ] **创建 Rule Response DTOs** (`internal/interface/http/dto/response/rule_response.go`)
  - [ ] 定义 `RuleResponse`（ID, UserID, RuleType, Name, Config, IsEnabled, CreatedAt, UpdatedAt）
  - [ ] 实现转换方法

### 5.4 HTTP Handler 实现

- [ ] **实现 NotificationRuleHandler** (`internal/interface/http/handler/notification_rule_handler.go`)
  - [ ] 注入依赖（NotificationRuleService, Validator, Logger）
  - [ ] 实现 `Create()` - POST /api/v1/notification-rules
  - [ ] 实现 `Update()` - PUT /api/v1/notification-rules/:id
  - [ ] 实现 `Delete()` - DELETE /api/v1/notification-rules/:id
  - [ ] 实现 `GetByID()` - GET /api/v1/notification-rules/:id
  - [ ] 实现 `List()` - GET /api/v1/notification-rules
  - [ ] 实现 `Toggle()` - PUT /api/v1/notification-rules/:id/toggle
  - [ ] 添加 Swagger 注解
  - [ ] 添加错误处理
  - [ ] 确保用户只能操作自己的规则

### 5.5 路由注册

- [ ] **更新 NotificationRoute** (`internal/interface/http/route/notification_route.go`)
  - [ ] 注册提醒规则相关路由

---

## Phase 6: 核心调度系统 (3-4天)

### 6.1 调度器配置

- [ ] **实现 Gocron 调度器** (`internal/infrastructure/scheduler/scheduler.go`)
  - [ ] 创建 `NewScheduler()` 函数
  - [ ] 配置定时任务（每分钟执行）
  - [ ] 集成 fx.Lifecycle（OnStart 启动调度器，OnStop 停止调度器）
  - [ ] 添加日志记录
  - [ ] 处理 panic 恢复

- [ ] **创建 Scheduler Module** (`internal/infrastructure/scheduler/module.go`)
  - [ ] 注册 Scheduler Provider

### 6.2 核心调度服务实现

- [ ] **实现 LiveCheckService** (`internal/application/service/live_check_service.go`)
  - [ ] 注入依赖（所有 Repositories, StreamingProviderManager, NotificationProviderManager, NotificationRuleService, Logger）
  - [ ] **实现 `CheckAllStreamers()`**
    - [ ] 获取所有有关注者的主播列表
    - [ ] 按平台分组
    - [ ] 使用 errgroup 并发检查各平台
    - [ ] 调用 `CheckStreamersByPlatform()`
    - [ ] 添加全局错误处理
  - [ ] **实现 `CheckStreamersByPlatform(platformId)`**
    - [ ] 获取该平台下的所有主播
    - [ ] 调用 Provider 的 `BatchCheckLiveStatus()`
    - [ ] 比对状态变化（离线→在线）
    - [ ] 对于开播事件，调用 `ProcessLiveOnEvent()`
    - [ ] 更新主播的直播状态和最后检查时间
  - [ ] **实现 `ProcessLiveOnEvent(streamer, liveStatus)`**
    - [ ] 获取该主播的所有关注者（notification_enabled=true）
    - [ ] 使用 errgroup 并发处理每个用户的通知
    - [ ] 调用 `SendNotificationToUser()`
  - [ ] **实现 `SendNotificationToUser(user, streamer, liveStatus)`**
    - [ ] 调用 `NotificationRuleService.ShouldNotify()` 判断
    - [ ] 如果不应该通知，直接返回
    - [ ] 获取用户的启用渠道（按优先级排序）
    - [ ] 依次尝试发送通知（失败则尝试下一个）
    - [ ] 成功后更新 `last_notified_at`
    - [ ] 记录发送结果日志
  - [ ] 添加错误处理和重试逻辑
  - [ ] 添加性能监控（执行时间、成功率）
  - [ ] 添加单元测试和集成测试

### 6.3 配置管理

- [ ] **更新配置文件** (`configs/config.yaml`)
  - [ ] 添加 `scheduler` 配置段
    - [ ] `live_check_interval`: 检查间隔（默认 60s）
  - [ ] 添加 `streaming` 配置段
    - [ ] `default_poll_interval`: 默认轮询间隔（默认 60秒）
    - [ ] `batch_size`: 批量检查主播数量（默认 50）
  - [ ] 添加 `notification` 配置段
    - [ ] `concurrent_send`: 并发发送通知数（默认 10）
    - [ ] `retry_count`: 重试次数（默认 3）

- [ ] **更新配置结构体** (`internal/infrastructure/config/config.go`)
  - [ ] 添加 SchedulerConfig 结构体
  - [ ] 添加 StreamingConfig 结构体
  - [ ] 添加 NotificationConfig 结构体

### 6.4 模块集成

- [ ] **更新主模块** (`internal/app/module.go`)
  - [ ] 导入 StreamingModule
  - [ ] 导入 NotificationModule
  - [ ] 导入 SchedulerModule
  - [ ] 确保依赖顺序正确

---

## Phase 7: 测试与优化 (2-3天)

### 7.1 单元测试

- [ ] **Domain 层测试**
  - [ ] 测试 Platform 实体
  - [ ] 测试 Streamer 实体
  - [ ] 测试 UserFollowing 实体
  - [ ] 测试 NotificationChannel 实体
  - [ ] 测试 NotificationRule 实体

- [ ] **Repository 层测试**
  - [ ] 使用 enttest 测试 PlatformRepository
  - [ ] 测试 StreamerRepository
  - [ ] 测试 FollowingRepository
  - [ ] 测试 ChannelRepository
  - [ ] 测试 RuleRepository

- [ ] **Service 层测试**
  - [ ] 测试 PlatformService
  - [ ] 测试 StreamerService
  - [ ] 测试 FollowingService
  - [ ] 测试 NotificationChannelService
  - [ ] 测试 NotificationRuleService
  - [ ] 测试 LiveCheckService（使用 mock）

- [ ] **Provider 测试**
  - [ ] 测试各个 StreamingPlatformProvider（使用 mock HTTP）
  - [ ] 测试各个 NotificationChannelProvider（使用 mock）

### 7.2 集成测试

- [ ] **API 集成测试**
  - [ ] 测试 Platform API endpoints
  - [ ] 测试 Streamer API endpoints
  - [ ] 测试 Following API endpoints
  - [ ] 测试 NotificationChannel API endpoints
  - [ ] 测试 NotificationRule API endpoints

- [ ] **端到端测试**
  - [ ] 测试完整的关注流程
  - [ ] 测试完整的通知流程
  - [ ] 测试规则过滤逻辑

### 7.3 性能优化

- [ ] **数据库优化**
  - [ ] 分析慢查询
  - [ ] 优化索引
  - [ ] 添加数据库连接池监控

- [ ] **并发优化**
  - [ ] 使用 errgroup 优化并发检查
  - [ ] 控制并发数量（避免过载）
  - [ ] 添加超时控制

- [ ] **缓存优化**（可选）
  - [ ] 考虑缓存主播信息
  - [ ] 考虑缓存用户关注列表
  - [ ] 使用 Redis（如果需要）

### 7.4 监控与日志

- [ ] **日志完善**
  - [ ] 添加结构化日志
  - [ ] 统一日志格式
  - [ ] 添加 trace_id 追踪

- [ ] **监控指标**
  - [ ] 添加调度任务执行时间监控
  - [ ] 添加通知发送成功率监控
  - [ ] 添加 API 响应时间监控
  - [ ] 添加错误率监控

- [ ] **告警配置**（可选）
  - [ ] 配置调度任务失败告警
  - [ ] 配置通知发送失败告警

### 7.5 错误处理优化

- [ ] **统一错误处理**
  - [ ] 检查所有错误处理是否统一
  - [ ] 确保错误信息对用户友好
  - [ ] 添加错误码文档

- [ ] **边界情况处理**
  - [ ] 处理网络超时
  - [ ] 处理并发冲突
  - [ ] 处理配置错误

### 7.6 文档完善

- [ ] **API 文档**
  - [ ] 生成 Swagger 文档：`make generate-swagger`
  - [ ] 验证所有端点都有文档
  - [ ] 添加 API 使用示例

- [ ] **配置文档**
  - [ ] 编写配置说明文档
  - [ ] 添加各个 Provider 的配置示例
  - [ ] 添加规则配置示例

- [ ] **部署文档**
  - [ ] 编写部署指南
  - [ ] 添加环境变量说明
  - [ ] 添加数据库迁移说明

---

## 额外任务（可选）

### 高级功能

- [ ] **通知历史记录**
  - [ ] 创建 NotificationHistory 实体
  - [ ] 记录所有发送的通知
  - [ ] 提供通知历史查询 API

- [ ] **用户偏好设置**
  - [ ] 支持全局静默时段
  - [ ] 支持通知摘要模式（一次推送多个开播）
  - [ ] 支持按平台或主播类型的偏好设置

- [ ] **管理后台**
  - [ ] 监控调度任务状态
  - [ ] 查看系统运行指标
  - [ ] 管理平台配置

### 性能优化

- [ ] **分布式调度**（大规模场景）
  - [ ] 使用分布式锁避免重复调度
  - [ ] 支持多实例部署

- [ ] **消息队列**（异步解耦）
  - [ ] 使用消息队列解耦检查和通知
  - [ ] 提高系统可靠性

---

## 验收标准

### 功能验收

- [ ] 用户可以关注多个平台的主播
- [ ] 用户可以配置多个推送渠道
- [ ] 用户可以设置提醒规则（静默时段、频率限制、内容过滤）
- [ ] 主播开播时，系统自动检测并推送通知
- [ ] 通知发送失败时自动切换到备用渠道
- [ ] 支持添加新的直播平台（3步扩展）
- [ ] 支持添加新的推送渠道（3步扩展）

### 性能验收

- [ ] 调度任务执行时间 < 30秒（500个主播）
- [ ] API 响应时间 < 200ms（P95）
- [ ] 通知发送成功率 > 99%
- [ ] 支持 1000+ 并发用户

### 质量验收

- [ ] 单元测试覆盖率 > 70%
- [ ] 所有 API 都有 Swagger 文档
- [ ] 所有错误都有友好的错误信息
- [ ] 日志完整且结构化

---

## 里程碑

| Phase | 目标 | 预计完成日期 | 实际完成日期 |
|-------|------|-------------|-------------|
| Phase 1 | 基础设施与数据模型 | ___ | ___ |
| Phase 2 | 直播平台集成 | ___ | ___ |
| Phase 3 | 关注功能 | ___ | ___ |
| Phase 4 | 推送渠道集成 | ___ | ___ |
| Phase 5 | 提醒规则 | ___ | ___ |
| Phase 6 | 核心调度系统 | ___ | ___ |
| Phase 7 | 测试与优化 | ___ | ___ |

---

## 依赖安装

```bash
# 安装新依赖
go get github.com/go-co-op/gocron/v2
go get github.com/go-resty/resty/v2
go get gopkg.in/gomail.v2
go get github.com/go-telegram-bot-api/telegram-bot-api/v5
go get golang.org/x/sync

# 更新依赖
go mod tidy
```

---

## 注意事项

1. **严格遵循 Clean Architecture**：确保依赖方向正确（Domain ← Application ← Infrastructure ← Interface）
2. **使用 fx 依赖注入**：所有组件通过 fx 自动装配，避免手动实例化
3. **错误处理**：使用项目统一的错误处理机制（`internal/pkg/errors/`）
4. **日志记录**：使用注入的 zap.Logger，记录关键操作和错误
5. **数据库事务**：需要事务的操作要正确使用 Ent 的事务机制
6. **测试先行**：关键业务逻辑先写测试再实现
7. **代码审查**：每个 Phase 完成后进行代码审查
8. **文档同步**：代码变更时同步更新文档

---

**祝开发顺利！** 🚀