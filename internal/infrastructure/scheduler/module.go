package scheduler

import (
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"go.uber.org/fx"
)

var Module = fx.Module("scheduler",
	fx.Provide(
		NewCronScheduler,

		fx.Annotate(
			NewJobRegistry,
			fx.ParamTags(`group:"jobs"`),
		),
	),

	fx.Invoke(func(cfg *config.Config, jobRegistry *JobRegistry) error {
		if err := jobRegistry.RegisterAll(cfg.Job); err != nil {
			return err
		}
		jobRegistry.scheduler.Start()
		return nil
	}),
)
