package command

type CreateNotificationChannelCommand struct {
	UserID      int64
	ChannelType string
	Name        string
	Config      map[string]any
	Enable      bool
	Priority    int
}

type UpdateNotificationChannelCommand struct {
	*CreateNotificationChannelCommand

	ID int64
}
