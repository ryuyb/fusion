package notification

import (
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/fx"
)

// Module provides the notification infrastructure components
var Module = fx.Module("notification",
	fx.Provide(
		// Notification channel providers
		// Each provider is annotated to be part of the "notification_providers" group
		fx.Annotate(
			NewWebhookProvider,
			fx.As(new(service.NotificationChannelProvider)),
			fx.ResultTags(`group:"notification_providers"`),
		),
	),
)
