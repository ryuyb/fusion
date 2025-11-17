package command

type CreateStreamingPlatformCommand struct {
	Type        string
	Name        string
	Description string
	BaseURL     string
	LogoURL     string
	Enabled     bool
	Priority    int
	Metadata    map[string]string
}

type UpdateStreamingPlatformCommand struct {
	*CreateStreamingPlatformCommand

	ID int64
}
