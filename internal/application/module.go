package application

import (
	"github.com/ryuyb/fusion/internal/application/job"
	"github.com/ryuyb/fusion/internal/application/service"
	"go.uber.org/fx"
)

var Module = fx.Module("application",
	fx.Provide(
		service.NewUserService,
		service.NewAuthService,
	),

	fx.Provide(
		asJob(job.NewBroadcastReminder),
	),
)

func asJob(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(job.Job)),
		fx.ResultTags(`group:"jobs"`),
	)
}
