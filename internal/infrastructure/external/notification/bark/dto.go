package bark

// BarkRequest
// refer: https://bark.day.app/#/tutorial?id=%e8%af%b7%e6%b1%82%e5%8f%82%e6%95%b0
type BarkRequest struct {
	Title     string `json:"title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
	Body      string `json:"body,omitempty"`
	DeviceKey string `json:"device_key,omitempty"`
	Level     string `json:"level,omitempty"`  // critical, active (default), timeSensitive, passive
	Volume    int    `json:"volume,omitempty"` // 0-10, default to 5
	Badge     int    `json:"badge,omitempty"`
	Call      string `json:"call,omitempty"`
	AutoCopy  string `json:"autoCopy,omitempty"`
	Copy      string `json:"copy,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	Url       string `json:"url,omitempty"`
	Action    string `json:"action,omitempty"`
}
