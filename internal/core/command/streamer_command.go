package command

type CreateStreamerCommand struct {
	PlatformType       string
	PlatformStreamerID string
	DisplayName        string
	AvatarURL          string
	RoomURL            string
	Bio                string
	Tags               []string
}

type UpdateStreamerCommand struct {
	*CreateStreamerCommand

	ID int64
}
