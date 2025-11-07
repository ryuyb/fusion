package client

import (
	"go.uber.org/fx"
)

var Module = fx.Module("client",
	fx.Provide(
		// Resty HTTP client for making API requests
		NewRestyClient,
	),
)
