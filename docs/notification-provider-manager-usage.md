# NotificationProviderManager 使用指南

## 概述

`NotificationProviderManager` 是通知系统中的核心组件，负责管理所有通知渠道提供者（Provider）。它通过 `fx` 依赖注入自动注册和管理所有的 `NotificationChannelProvider` 实现。

## 功能特性

- **自动注册**：通过 `fx.Group("notification_providers")` 自动收集所有 Provider
- **类型安全**：使用 `ChannelType` 枚举确保类型安全
- **灵活查询**：支持按类型获取 Provider 或获取所有 Provider
- **日志集成**：集成 zap 日志记录所有操作
- **可扩展性**：新增 Provider 时无需修改 Manager 代码

## 核心方法

### 1. GetProvider(channelType)
根据渠道类型获取对应的 Provider。

```go
provider, err := manager.GetProvider(entity.ChannelTypeWebhook)
if err != nil {
    return err
}
```

### 2. GetAllProviders()
获取所有已注册的 Provider 列表。

```go
providers := manager.GetAllProviders()
for _, provider := range providers {
    channelType := provider.GetChannelType()
    // 处理每个 provider
}
```

### 3. HasProvider(channelType)
检查指定类型的 Provider 是否已注册。

```go
if manager.HasProvider(entity.ChannelTypeTelegram) {
    // 有 Telegram Provider
}
```

### 4. GetSupportedChannels()
获取所有支持的渠道类型列表。

```go
supportedChannels := manager.GetSupportedChannels()
for _, channelType := range supportedChannels {
    log.Printf("支持渠道: %s", channelType)
}
```

## 使用场景

### 场景1：在应用服务中使用

```go
type NotificationChannelService struct {
    manager     *notification.NotificationProviderManager
    repo        repository.ChannelRepository
    logger      *zap.Logger
}

func (s *NotificationChannelService) CreateChannel(userID int64, req *CreateChannelRequest) (*entity.NotificationChannel, error) {
    // 验证 Provider 是否存在
    if !s.manager.HasProvider(req.ChannelType) {
        return nil, fmt.Errorf("不支持的渠道类型: %s", req.ChannelType)
    }

    // 获取 Provider 并验证配置
    provider, err := s.manager.GetProvider(req.ChannelType)
    if err != nil {
        return nil, err
    }

    if err := provider.ValidateConfiguration(req.Config); err != nil {
        return nil, fmt.Errorf("配置验证失败: %w", err)
    }

    // 创建渠道
    channel := entity.CreateChannel(userID, req.ChannelType, req.Name, req.Config, req.Priority)
    // ... 保存到数据库
    return channel, nil
}
```

### 场景2：发送通知

```go
func (s *NotificationChannelService) SendNotification(userID int64, notification *service.Notification) error {
    // 获取用户启用的渠道
    channels, err := s.repo.FindEnabledByUser(ctx, userID)
    if err != nil {
        return err
    }

    // 按优先级排序
    sort.Slice(channels, func(i, j int) bool {
        return channels[i].Priority < channels[j].Priority
    })

    // 依次尝试发送
    var lastErr error
    for _, channel := range channels {
        provider, err := s.manager.GetProvider(channel.ChannelType)
        if err != nil {
            lastErr = err
            continue
        }

        if err := provider.Send(ctx, channel, notification); err != nil {
            s.logger.Warn("发送通知失败", zap.String("channel", channel.Name), zap.Error(err))
            lastErr = err
            continue
        }

        s.logger.Info("通知发送成功", zap.String("channel", channel.Name))
        return nil // 成功发送，退出
    }

    return fmt.Errorf("所有渠道发送失败，最后错误: %w", lastErr)
}
```

### 场景3：测试渠道连通性

```go
func (s *NotificationChannelService) TestChannel(userID, channelID int64) error {
    channel, err := s.repo.FindByID(ctx, channelID)
    if err != nil {
        return err
    }

    // 确保用户只能测试自己的渠道
    if channel.UserID != userID {
        return errors.Unauthorized("无权操作此渠道")
    }

    provider, err := s.manager.GetProvider(channel.ChannelType)
    if err != nil {
        return err
    }

    return provider.TestConnection(ctx, channel.Config)
}
```

## Provider 注册

### 现有 Provider

1. **WebhookProvider** - Webhook 推送
   - 类型：`entity.ChannelTypeWebhook`
   - 配置：URL, HTTP方法, Headers, 模板等

### 未来扩展的 Provider

1. **EmailProvider** - 邮件推送
   - 类型：`entity.ChannelTypeEmail`
   - 配置：SMTP服务器、用户名、密码等

2. **TelegramProvider** - Telegram Bot
   - 类型：`entity.ChannelTypeTelegram`
   - 配置：Bot Token, Chat ID 等

3. **DiscordProvider** - Discord Webhook
   - 类型：`entity.ChannelTypeDiscord`
   - 配置：Webhook URL 等

4. **FeishuProvider** - 飞书机器人
   - 类型：`entity.ChannelTypeFeishu`
   - 配置：Webhook URL 等

## 添加新 Provider 的步骤

1. **创建 Provider 实现**
   ```go
   type MyChannelProvider struct {
       client *client.RestyClient
       logger *zap.Logger
   }

   func (p *MyChannelProvider) GetChannelType() entity.ChannelType {
       return entity.ChannelTypeMyCustom
   }

   // 实现其他接口方法...
   ```

2. **在 module.go 中注册**
   ```go
   var Module = fx.Module("notification",
       fx.Provide(
           fx.Annotate(
               NewMyChannelProvider,
               fx.As(new(service.NotificationChannelProvider)),
               fx.ResultTags(`group:"notification_providers"`),
           ),
           // ... 其他 providers
           NewNotificationProviderManager,
       ),
   )
   ```

3. **无需修改 Manager**
   Manager 会自动通过 `fx.Group` 收集并注册新 Provider。

## 最佳实践

1. **错误处理**：始终检查 Provider 是否存在再使用
2. **配置验证**：在创建渠道前先验证配置
3. **优先级排序**：发送通知时按优先级排序渠道
4. **失败重试**：一个渠道失败时尝试下一个渠道
5. **日志记录**：记录所有关键操作和错误
6. **测试覆盖**：为 Provider 和 Manager 编写测试

## 与 StreamingProviderManager 对比

| 特性 | NotificationProviderManager | StreamingProviderManager |
|------|----------------------------|--------------------------|
| 管理的类型 | 通知渠道 (Email, Webhook 等) | 直播平台 (Douyu, Huya 等) |
| 接口 | NotificationChannelProvider | StreamingPlatformProvider |
| 枚举类型 | ChannelType | PlatformType |
| 主要用途 | 发送开播通知 | 获取直播状态 |

两者设计模式完全一致，都使用 `fx.Group` 自动注册，便于扩展和维护。