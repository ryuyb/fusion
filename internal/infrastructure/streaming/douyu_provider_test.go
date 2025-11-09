package streaming

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/infrastructure/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func setupTestDouyuProvider(t *testing.T) *DouyuProvider {
	logger := zaptest.NewLogger(t)
	restyClient := client.NewRestyClient(logger)
	return NewDouyuProvider(restyClient, logger)
}

func TestFetchStreamerInfo(t *testing.T) {
	provider := setupTestDouyuProvider(t)

	info, err := provider.FetchStreamerInfo(context.Background(), "60937")

	require.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, info.PlatformStreamerId, "60937")
	assert.Equal(t, info.Name, "zard1991")
	assert.NotEmpty(t, info.Avatar)
	assert.NotEmpty(t, info.Description)
	assert.Equal(t, info.RoomURL, "https://www.douyu.com/60937")
}

func TestCheckLiveStatus(t *testing.T) {
	provider := setupTestDouyuProvider(t)

	liveStatus, err := provider.CheckLiveStatus(context.Background(), "60937")

	require.NoError(t, err)
	assert.NotNil(t, liveStatus)
}
