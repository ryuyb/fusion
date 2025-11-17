package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserFollowedStreamer(t *testing.T) {
	follow, err := NewUserFollowedStreamer(1, 2, " alias ", " notes ", []int64{10, 10, 20})
	require.NoError(t, err)

	require.Equal(t, int64(1), follow.UserID)
	require.Equal(t, int64(2), follow.StreamerID)
	require.Equal(t, "alias", follow.Alias)
	require.Equal(t, "notes", follow.Notes)
	require.Equal(t, []int64{10, 20}, follow.NotificationChannelIDs)
	require.True(t, follow.NotificationsEnabled)
}

func TestNewUserFollowedStreamerValidation(t *testing.T) {
	cases := []struct {
		name        string
		userID      int64
		streamerID  int64
		channelIDs  []int64
		expectedMsg string
	}{
		{
			name:        "invalid user id",
			userID:      0,
			streamerID:  1,
			expectedMsg: "user id must be greater than zero",
		},
		{
			name:        "invalid streamer id",
			userID:      1,
			streamerID:  0,
			expectedMsg: "streamer id must be greater than zero",
		},
		{
			name:        "invalid notification channel id",
			userID:      1,
			streamerID:  2,
			channelIDs:  []int64{-1},
			expectedMsg: "notification channel id must be greater than zero",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUserFollowedStreamer(tc.userID, tc.streamerID, "", "", tc.channelIDs)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedMsg)
		})
	}
}

func TestUserFollowedStreamerUpdatePreferences(t *testing.T) {
	follow, err := NewUserFollowedStreamer(1, 2, "alias", "notes", []int64{1})
	require.NoError(t, err)

	err = follow.UpdatePreferences("  alias-2  ", "  notes-2  ", false, []int64{5, 5, 7})
	require.NoError(t, err)

	require.Equal(t, "alias-2", follow.Alias)
	require.Equal(t, "notes-2", follow.Notes)
	require.False(t, follow.NotificationsEnabled)
	require.Equal(t, []int64{5, 7}, follow.NotificationChannelIDs)

	require.Error(t, follow.UpdatePreferences("", "", true, []int64{0}))
}
