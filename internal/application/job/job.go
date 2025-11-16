package job

import "context"

type Job interface {
	Name() string

	Execute(ctx context.Context) error
}
