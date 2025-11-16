package notification

import (
	"github.com/ryuyb/fusion/internal/core/port/external"
	"github.com/ryuyb/fusion/internal/infrastructure/external/notification/bark"
	"go.uber.org/fx"
)

var Module = fx.Module("notification",
	fx.Provide(
		asProvider(bark.NewProvider),
	),

	fx.Provide(
		fx.Annotate(
			NewNotificationProviderManager,
			fx.ParamTags(`group:"notification_providers"`),
		),
	),
)

func asProvider(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(external.NotificationProvider)),
		fx.ResultTags(`group:"notification_providers"`),
	)
}
