package client

import (
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// RestyClient wraps resty.Client for streaming platform API calls
type RestyClient struct {
	*resty.Client
}

// NewRestyClient creates a new configured Resty HTTP client
func NewRestyClient(logger *zap.Logger) *RestyClient {
	client := resty.New()

	// Set timeout
	client.SetTimeout(10 * time.Second)

	// Set retry configuration
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)

	// Set User-Agent header
	client.SetHeader("User-Agent", "Fusion-Streaming-Platform/1.0")

	// Add request logging middleware
	client.AddRequestMiddleware(func(client *resty.Client, req *resty.Request) error {
		logger.Debug("Outgoing HTTP req",
			zap.String("method", req.Method),
			zap.String("url", req.URL),
			zap.String("params", req.QueryParams.Encode()),
		)
		return nil
	})

	// Add response logging middleware
	client.AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
		logger.Debug("Incoming HTTP response",
			zap.String("method", resp.Request.Method),
			zap.String("url", resp.Request.URL),
			zap.Int("status_code", resp.StatusCode()),
			zap.Duration("time", resp.Duration()),
		)
		return nil
	})

	// Add error logging middleware
	client.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logger.Error("HTTP request failed",
				zap.String("method", req.Method),
				zap.String("url", req.URL),
				zap.Int("status_code", v.Response.StatusCode()),
				zap.Error(err),
			)
		} else {
			logger.Error("HTTP request error",
				zap.String("method", req.Method),
				zap.String("url", req.URL),
				zap.Error(err),
			)
		}
	})

	return &RestyClient{Client: client}
}
