package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/ryuyb/fusion/internal/pkg/entgo/mixin"
)

// NotificationRule holds the schema definition for the NotificationRule entity.
type NotificationRule struct {
	ent.Schema
}

// Fields of the NotificationRule.
func (NotificationRule) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.Int64("user_id").
			Positive(),
		field.Enum("rule_type").
			Values("silent_period", "rate_limit", "content_filter"),
		field.String("name").
			NotEmpty(),
		field.JSON("config", map[string]interface{}{}).
			Optional(),
		field.Bool("is_enabled").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the NotificationRule.
func (NotificationRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "is_enabled"),
		index.Fields("rule_type"),
	}
}

// Edges of the NotificationRule.
func (NotificationRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("notification_rules").
			Field("user_id").
			Required().
			Unique(),
	}
}

// Mixin of the NotificationRule.
func (NotificationRule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}
