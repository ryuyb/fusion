package streaming

import (
	"github.com/ryuyb/fusion/internal/core/port/external"
	"github.com/ryuyb/fusion/internal/infrastructure/external/streaming/bilibili"
	"github.com/ryuyb/fusion/internal/infrastructure/external/streaming/douyu"
	"go.uber.org/fx"
)

var Module = fx.Module("streaming",
	fx.Provide(
		asProvider(bilibili.NewProvider),
		asProvider(douyu.NewProvider),
	),

	fx.Provide(
		fx.Annotate(
			NewStreamingProviderManager,
			fx.ParamTags(`group:"streaming_providers"`),
		),
	),
)

func asProvider(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(external.StreamingPlatformProvider)),
		fx.ResultTags(`group:"streaming_providers"`),
	)
}
