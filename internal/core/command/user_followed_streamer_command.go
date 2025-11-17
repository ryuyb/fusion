package command

type CreateUserFollowedStreamerCommand struct {
	UserID                 int64
	StreamerID             int64
	Alias                  string
	Notes                  string
	NotificationsEnabled   bool
	NotificationChannelIDs []int64
}

type UpdateUserFollowedStreamerCommand struct {
	ID                     int64
	Alias                  string
	Notes                  string
	NotificationsEnabled   bool
	NotificationChannelIDs []int64
}
