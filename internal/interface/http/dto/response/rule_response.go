package response

import (
	"time"

	"github.com/ryuyb/fusion/internal/domain/entity"
)

// RuleResponse represents a notification rule in the response
type RuleResponse struct {
	ID        int64                  `json:"id"`
	UserID    int64                  `json:"user_id"`
	RuleType  string                 `json:"rule_type"`
	Name      string                 `json:"name"`
	Config    map[string]interface{} `json:"config"`
	IsEnabled bool                   `json:"is_enabled"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// RuleConfigInfo represents detailed information about a rule type
type RuleConfigInfo struct {
	RuleType    string                 `json:"rule_type"`
	TypeName    string                 `json:"rule_type_name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Fields      []RuleConfigField      `json:"fields"`
}

type RuleConfigField struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

// ToRuleResponse converts an entity.NotificationRule to RuleResponse
func ToRuleResponse(rule *entity.NotificationRule) *RuleResponse {
	return &RuleResponse{
		ID:        rule.ID,
		UserID:    rule.UserID,
		RuleType:  string(rule.RuleType),
		Name:      rule.Name,
		Config:    rule.Config,
		IsEnabled: rule.IsEnabled,
		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,
	}
}

// ToRuleResponseList converts a list of entities to response with pagination
func ToRuleResponseList(rules []*entity.NotificationRule, total, page, pageSize int) *PaginationResponse[*RuleResponse] {
	responses := make([]*RuleResponse, 0, len(rules))
	for _, rule := range rules {
		responses = append(responses, ToRuleResponse(rule))
	}

	return NewPaginationResponse[*RuleResponse](responses, total, page, pageSize)
}

// ToRuleConfigInfo converts a rule entity to detailed config info
func ToRuleConfigInfo(rule *entity.NotificationRule) *RuleConfigInfo {
	info := &RuleConfigInfo{
		RuleType: string(rule.RuleType),
		Config:   rule.Config,
		Fields:   getFieldsForRuleType(rule.RuleType),
	}

	// Set type name and description based on rule type
	switch rule.RuleType {
	case entity.RuleTypeSilentPeriod:
		info.TypeName = "静默时段"
		info.Description = "设置不发送通知的时间段，例如夜间23:00到早上08:00"
	case entity.RuleTypeRateLimit:
		info.TypeName = "频率限制"
		info.Description = "限制通知发送频率，避免过度通知"
	case entity.RuleTypeContentFilter:
		info.TypeName = "内容过滤"
		info.Description = "根据直播标题或分类过滤通知"
	}

	return info
}

// getFieldsForRuleType returns the field definitions for a specific rule type
func getFieldsForRuleType(ruleType entity.RuleType) []RuleConfigField {
	switch ruleType {
	case entity.RuleTypeSilentPeriod:
		return []RuleConfigField{
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
		}
	case entity.RuleTypeRateLimit:
		return []RuleConfigField{
			{
				Field:       "interval_hours",
				Type:        "int",
				Required:    true,
				Description: "通知间隔时间（小时）",
				Example:     "2",
			},
		}
	case entity.RuleTypeContentFilter:
		return []RuleConfigField{
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
		}
	default:
		return []RuleConfigField{}
	}
}
