package streaming

import (
	"context"
	"testing"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/infrastructure/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func setupTestBilibiliProvider(t *testing.T) *BilibiliProvider {
	logger := zaptest.NewLogger(t)
	client := client.NewRestyClient(logger)
	return NewBilibiliProvider(client, logger)
}

// TestBilibiliProvider_GetPlatformType tests that the provider returns correct platform type
func TestBilibiliProvider_GetPlatformType(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	platformType := provider.GetPlatformType()

	assert.Equal(t, entity.PlatformTypeBilibili, platformType)
}

// TestBilibiliProvider_ValidateConfiguration tests configuration validation
func TestBilibiliProvider_ValidateConfiguration(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	tests := []struct {
		name   string
		config map[string]interface{}
	}{
		{
			name:   "nil config",
			config: nil,
		},
		{
			name:   "empty config",
			config: map[string]interface{}{},
		},
		{
			name: "any config",
			config: map[string]interface{}{
				"some_key": "some_value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.ValidateConfiguration(tt.config)
			require.NoError(t, err, "Bilibili doesn't require configuration validation")
		})
	}
}

// TestBilibiliProvider_SearchStreamer tests that search is not supported
func TestBilibiliProvider_SearchStreamer(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	ctx := context.Background()
	results, err := provider.SearchStreamer(ctx, "测试")

	require.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "not supported", "Should indicate that search is not supported")
}

// TestBilibiliProvider_FetchStreamerInfo_InvalidRoomID tests invalid room ID handling
func TestBilibiliProvider_FetchStreamerInfo_InvalidRoomID(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	tests := []struct {
		name   string
		roomID string
	}{
		{
			name:   "non-numeric room ID",
			roomID: "invalid",
		},
		{
			name:   "empty room ID",
			roomID: "",
		},
		{
			name:   "special characters",
			roomID: "abc@123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := provider.FetchStreamerInfo(ctx, tt.roomID)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid room id")
		})
	}
}

// TestBilibiliProvider_CheckLiveStatus_InvalidRoomID tests invalid room ID handling
func TestBilibiliProvider_CheckLiveStatus_InvalidRoomID(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	tests := []struct {
		name        string
		roomID      string
		errorSubstr string
	}{
		{
			name:        "non-numeric room ID",
			roomID:      "invalid",
			errorSubstr: "invalid room id",
		},
		{
			name:        "empty room ID",
			roomID:      "",
			errorSubstr: "invalid room id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := provider.CheckLiveStatus(ctx, tt.roomID)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorSubstr)
		})
	}
}

// TestBilibiliProvider_BatchCheckLiveStatus_EmptyList tests empty input handling
func TestBilibiliProvider_BatchCheckLiveStatus_EmptyList(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	ctx := context.Background()
	results, err := provider.BatchCheckLiveStatus(ctx, []string{})

	require.NoError(t, err)
	assert.Empty(t, results)
}

// TestBilibiliProvider_BatchCheckLiveStatus_NilList tests nil input handling
func TestBilibiliProvider_BatchCheckLiveStatus_NilList(t *testing.T) {
	provider := setupTestBilibiliProvider(t)

	ctx := context.Background()
	results, err := provider.BatchCheckLiveStatus(ctx, nil)

	require.NoError(t, err)
	assert.Empty(t, results)
}

// Integration tests - these require network access to Bilibili API
// Run with: go test -v -tags=integration

// TestBilibiliProvider_FetchStreamerInfo_Integration tests real API integration
// Using a well-known streamer room ID: 7777 (笑笑)
func TestBilibiliProvider_FetchStreamerInfo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := setupTestBilibiliProvider(t)
	ctx := context.Background()

	// Test with a well-known Bilibili room
	roomID := "7777"
	info, err := provider.FetchStreamerInfo(ctx, roomID)

	require.NoError(t, err, "Should successfully fetch streamer info")
	assert.NotNil(t, info)
	assert.Equal(t, roomID, info.PlatformStreamerId)
	assert.NotEmpty(t, info.Name, "Streamer name should not be empty")
	assert.NotEmpty(t, info.RoomURL, "Room URL should not be empty")
	assert.Contains(t, info.RoomURL, roomID, "Room URL should contain room ID")

	t.Logf("Fetched streamer info: Name=%s, RoomURL=%s", info.Name, info.RoomURL)
}

// TestBilibiliProvider_CheckLiveStatus_Integration tests live status checking
func TestBilibiliProvider_CheckLiveStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := setupTestBilibiliProvider(t)
	ctx := context.Background()

	// Test with a well-known Bilibili room
	roomID := "7777"
	status, err := provider.CheckLiveStatus(ctx, roomID)

	require.NoError(t, err, "Should successfully check live status")
	assert.NotNil(t, status)

	// IsLive can be true or false depending on current status
	t.Logf("Live status: IsLive=%v, Title=%s, Viewers=%d",
		status.IsLive, status.Title, status.Viewers)

	// Basic validations
	assert.NotEmpty(t, status.Title, "Title should not be empty")
	if status.IsLive {
		assert.GreaterOrEqual(t, status.Viewers, 0, "Viewers count should be non-negative")
		assert.False(t, status.StartTime.IsZero(), "Start time should be set when live")
	}
}

// TestBilibiliProvider_CheckLiveStatus_NonExistentRoom tests error handling for non-existent room
func TestBilibiliProvider_CheckLiveStatus_NonExistentRoom(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := setupTestBilibiliProvider(t)
	ctx := context.Background()

	// Use a very unlikely to exist room ID
	roomID := "9999999999"
	status, err := provider.CheckLiveStatus(ctx, roomID)

	// Bilibili API might return an error or return data with error code
	if err != nil {
		assert.Contains(t, err.Error(), "API error", "Should indicate API error")
		t.Logf("Received expected error: %v", err)
	} else {
		// If no error, the status might have default values
		assert.NotNil(t, status)
		t.Logf("Received status for non-existent room: %+v", status)
	}
}

// TestBilibiliProvider_BatchCheckLiveStatus_Integration tests batch checking
func TestBilibiliProvider_BatchCheckLiveStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := setupTestBilibiliProvider(t)
	ctx := context.Background()

	// Test with multiple well-known rooms
	roomIDs := []string{"7777", "5096"}
	results, err := provider.BatchCheckLiveStatus(ctx, roomIDs)

	require.NoError(t, err, "Should successfully batch check live status")

	// We expect results for all rooms (even if some fail, BatchCheckLiveStatus continues)
	t.Logf("Got results for %d out of %d rooms", len(results), len(roomIDs))

	// Check each result
	for roomID, status := range results {
		assert.NotNil(t, status, "Status should not be nil for room %s", roomID)
		t.Logf("Room %s: IsLive=%v, Title=%s", roomID, status.IsLive, status.Title)
	}
}

// TestBilibiliProvider_BatchCheckLiveStatus_MixedValidity tests batch with invalid IDs
func TestBilibiliProvider_BatchCheckLiveStatus_MixedValidity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := setupTestBilibiliProvider(t)
	ctx := context.Background()

	// Mix of valid and potentially invalid room IDs
	roomIDs := []string{"7777", "9999999999", "5096"}
	results, err := provider.BatchCheckLiveStatus(ctx, roomIDs)

	// Should not error out completely - it continues on failures
	require.NoError(t, err)

	// Should have at least some results
	assert.NotEmpty(t, results, "Should have results for valid rooms")
	t.Logf("Got results for %d out of %d rooms", len(results), len(roomIDs))

	// Check that valid rooms have results
	for roomID, status := range results {
		assert.NotNil(t, status, "Status should not be nil for room %s", roomID)
	}
}

// Benchmark tests

func BenchmarkBilibiliProvider_CheckLiveStatus(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	provider := setupTestBilibiliProvider(&testing.T{})
	ctx := context.Background()
	roomID := "7777"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.CheckLiveStatus(ctx, roomID)
		if err != nil {
			b.Logf("Error checking status: %v", err)
		}
	}
}

func BenchmarkBilibiliProvider_BatchCheckLiveStatus(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	provider := setupTestBilibiliProvider(&testing.T{})
	ctx := context.Background()
	roomIDs := []string{"7777", "5096", "21852"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.BatchCheckLiveStatus(ctx, roomIDs)
		if err != nil {
			b.Logf("Error batch checking: %v", err)
		}
	}
}
