package request

// CreateRuleRequest represents the request to create a notification rule
type CreateRuleRequest struct {
	RuleType string                 `json:"rule_type" validate:"required,oneof=silent_period rate_limit content_filter"`
	Name     string                 `json:"name" validate:"required,min=2,max=50"`
	Config   map[string]interface{} `json:"config" validate:"required"`
}

// UpdateRuleRequest represents the request to update a notification rule
type UpdateRuleRequest struct {
	Name      string                 `json:"name" validate:"required,min=2,max=50"`
	Config    map[string]interface{} `json:"config" validate:"required"`
	IsEnabled *bool                  `json:"is_enabled" validate:"omitempty"`
}

// ToggleRuleRequest represents the request to toggle a rule enabled status
type ToggleRuleRequest struct {
	IsEnabled bool `json:"is_enabled" validate:"required"`
}

// ListRuleRequest represents the request to list notification rules
type ListRuleRequest struct {
	RuleType  string `json:"rule_type" validate:"omitempty,oneof=silent_period rate_limit content_filter"`
	IsEnabled *bool  `json:"is_enabled" validate:"omitempty"`
	Page      int    `json:"page" validate:"omitempty,gte=1"`
	PageSize  int    `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}

// RuleConfigExample provides examples for different rule types
type RuleConfigExample struct {
	RuleType      string                 `json:"rule_type"`
	RuleTypeName  string                 `json:"rule_type_name"`
	Description   string                 `json:"description"`
	ExampleConfig map[string]interface{} `json:"example_config"`
	Fields        []RuleConfigField      `json:"fields"`
}

type RuleConfigField struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

// GetConfigExamples returns configuration examples for all rule types
func GetConfigExamples() []RuleConfigExample {
	return []RuleConfigExample{
		{
			RuleType:     "silent_period",
			RuleTypeName: "静默时段",
			Description:  "设置不发送通知的时间段，例如夜间23:00到早上08:00",
			ExampleConfig: map[string]interface{}{
				"start_hour": 23,
				"end_hour":   8,
			},
			Fields: []RuleConfigField{
				{
					Field:       "start_hour",
					Type:        "int",
					Required:    true,
					Description: "静默时段的开始时间（0-23）",
					Example:     "23",
				},
				{
					Field:       "end_hour",
					Type:        "int",
					Required:    true,
					Description: "静默时段的结束时间（0-23）",
					Example:     "8",
				},
			},
		},
		{
			RuleType:     "rate_limit",
			RuleTypeName: "频率限制",
			Description:  "限制通知发送频率，避免过度通知",
			ExampleConfig: map[string]interface{}{
				"interval_hours": 2,
			},
			Fields: []RuleConfigField{
				{
					Field:       "interval_hours",
					Type:        "int",
					Required:    true,
					Description: "通知间隔时间（小时）",
					Example:     "2",
				},
			},
		},
		{
			RuleType:     "content_filter",
			RuleTypeName: "内容过滤",
			Description:  "根据直播标题或分类过滤通知",
			ExampleConfig: map[string]interface{}{
				"keywords":   []string{"英雄联盟", "LOL", "王者"},
				"match_mode": "any",
			},
			Fields: []RuleConfigField{
				{
					Field:       "keywords",
					Type:        "[]string",
					Required:    true,
					Description: "匹配关键词列表",
					Example:     "[\"英雄联盟\", \"LOL\"]",
				},
				{
					Field:       "match_mode",
					Type:        "string",
					Required:    true,
					Description: "匹配模式：any(任一匹配) 或 all(全部匹配)",
					Example:     "any",
				},
			},
		},
	}
}
