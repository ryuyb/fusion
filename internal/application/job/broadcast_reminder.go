package job

import (
	"context"

	"go.uber.org/zap"
)

const (
	BroadcastReminderJob = "broadcast_reminder"
)

type BroadcastReminder struct {
	logger *zap.Logger
}

func NewBroadcastReminder(logger *zap.Logger) *BroadcastReminder {
	return &BroadcastReminder{
		logger: logger,
	}
}

func (j *BroadcastReminder) Name() string {
	return BroadcastReminderJob
}

func (j *BroadcastReminder) Execute(ctx context.Context) error {
	j.logger.Info("BroadcastReminderJob")
	return nil
}
