package scheduler

import (
	"fmt"

	"github.com/ryuyb/fusion/internal/application/job"
	"github.com/ryuyb/fusion/internal/infrastructure/provider/config"
	"go.uber.org/zap"
)

type JobRegistry struct {
	jobs      map[string]job.Job
	scheduler *CronScheduler
	logger    *zap.Logger
}

func NewJobRegistry(jobs []job.Job, scheduler *CronScheduler, logger *zap.Logger) *JobRegistry {
	jobsMap := make(map[string]job.Job, len(jobs))
	for _, item := range jobs {
		jobsMap[item.Name()] = item
	}
	return &JobRegistry{
		jobs:      jobsMap,
		scheduler: scheduler,
		logger:    logger,
	}
}

func (r *JobRegistry) RegisterAll(jobConfigs map[string]config.JobConfig) error {
	for name, jobConfig := range jobConfigs {
		if !jobConfig.Enable {
			r.logger.Info("Job disabled, skipping", zap.String("name", name))
			continue
		}
		jobInstance, ok := r.jobs[name]
		if !ok {
			return fmt.Errorf("unknown job type: %s", name)
		}
		if err := r.scheduler.Register(jobConfig.CronExpr, jobInstance); err != nil {
			return fmt.Errorf("register cron expression failed: %w", err)
		}
	}
	return nil
}
