package external

import (
	"github.com/ryuyb/fusion/internal/infrastructure/external/notification"
	"github.com/ryuyb/fusion/internal/infrastructure/external/streaming"
	"go.uber.org/fx"
)

var Module = fx.Module("external",
	notification.Module,
	streaming.Module,
)
