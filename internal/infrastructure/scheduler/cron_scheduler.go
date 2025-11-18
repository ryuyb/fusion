package scheduler

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/ryuyb/fusion/internal/application/job"
	"go.uber.org/zap"
)

type CronScheduler struct {
	scheduler gocron.Scheduler
	logger    *zap.Logger
}

func NewCronScheduler(logger *zap.Logger) (*CronScheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Error("failed to create scheduler", zap.Error(err))
		return nil, err
	}
	return &CronScheduler{
		scheduler: scheduler,
		logger:    logger,
	}, nil
}

func (s *CronScheduler) Register(cronExpr string, job job.Job) error {
	_, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false),
		gocron.NewTask(func(ctx context.Context) {
			s.executeJob(ctx, job)
		}),
		gocron.WithName(job.Name()),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		s.logger.Error("failed to register job", zap.String("name", job.Name()), zap.Error(err))
		return err
	}
	s.logger.Info("job registered", zap.String("name", job.Name()), zap.String("cron", cronExpr))
	return nil
}

func (s *CronScheduler) executeJob(ctx context.Context, job job.Job) {
	s.logger.Info("Job started", zap.String("name", job.Name()))

	err := job.Execute(ctx)

	if err != nil {
		s.logger.Error("Job failed", zap.String("name", job.Name()), zap.Error(err))
	} else {
		s.logger.Info("Job completed", zap.String("name", job.Name()))
	}
}

func (s *CronScheduler) Start() {
	s.logger.Info("Starting scheduler")
	s.scheduler.Start()
	s.logger.Info("Scheduler started")
}

func (s *CronScheduler) Stop() {
	s.logger.Info("Stopping scheduler")
	if err := s.scheduler.StopJobs(); err != nil {
		s.logger.Error("failed to stop scheduler", zap.Error(err))
	} else {
		s.logger.Info("Scheduler stopped")
	}
}

func (s *CronScheduler) Shutdown() {
	s.logger.Info("Shutting down scheduler")
	if err := s.scheduler.Shutdown(); err != nil {
		s.logger.Error("failed to shutdown scheduler", zap.Error(err))
	} else {
		s.logger.Info("Scheduler shutdown completed")
	}
}

func (s *CronScheduler) ListJobs() []gocron.Job {
	return s.scheduler.Jobs()
}
