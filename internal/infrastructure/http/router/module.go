package router

import "go.uber.org/fx"

var Module = fx.Module("router",
	fx.Provide(
		asRouter(NewHealthRouter),
		asRouter(NewUserRouter),
		asRouter(NewAuthRouter),
		asRouter(NewStreamingPlatformRouter),
		asRouter(NewStreamerRouter),
		asRouter(NewNotificationChannelRouter),
		asRouter(NewUserFollowedStreamerRouter),
	),

	fx.Provide(NewRouterRegistry),
)

func asRouter(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Router)),
		fx.ResultTags(`group:"routers"`),
	)
}
