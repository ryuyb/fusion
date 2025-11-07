package entity

import (
	"time"
)

// NotificationRule represents a user's notification rule configuration
type NotificationRule struct {
	ID        int64
	UserID    int64
	RuleType  RuleType
	Name      string
	Config    map[string]interface{} // JSON configuration for the rule
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeleteAt  time.Time
}

// RuleType defines the type of notification rule
type RuleType string

const (
	RuleTypeSilentPeriod  RuleType = "silent_period"  // Silent period rule (e.g., 23:00-08:00)
	RuleTypeRateLimit     RuleType = "rate_limit"     // Rate limit rule (e.g., max once per 2 hours)
	RuleTypeContentFilter RuleType = "content_filter" // Content filter rule (e.g., keywords matching)
)

// CreateRule creates a new NotificationRule instance
func CreateRule(userID int64, ruleType RuleType, name string, config map[string]interface{}) *NotificationRule {
	return &NotificationRule{
		UserID:    userID,
		RuleType:  ruleType,
		Name:      name,
		Config:    config,
		IsEnabled: true, // Default to enabled
	}
}

// Update updates the rule information
func (nr *NotificationRule) Update(name string, config map[string]interface{}) *NotificationRule {
	nr.Name = name
	nr.Config = config
	return nr
}

// Toggle toggles the enabled status of the rule
func (nr *NotificationRule) Toggle(enabled bool) *NotificationRule {
	nr.IsEnabled = enabled
	return nr
}

// IsActive checks if the rule is enabled
func (nr *NotificationRule) IsActive() bool {
	return nr.IsEnabled
}
