package streaming

import (
	"github.com/ryuyb/fusion/internal/domain/service"
	"go.uber.org/fx"
)

// Module provides the streaming infrastructure components
var Module = fx.Module("streaming",
	fx.Provide(
		// Streaming platform providers (will be implemented in Phase 2.2)
		// Each provider is annotated to be part of the "streaming_providers" group
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
		fx.Annotate(
			NewStreamingProviderManager,
			fx.ParamTags(`group:"streaming_providers"`),
		),
	),
)
