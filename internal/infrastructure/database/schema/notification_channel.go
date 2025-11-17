package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NotificationChannel stores delivery endpoints for a user's alerts.
type NotificationChannel struct {
	ent.Schema
}

func (NotificationChannel) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.Int64("user_id").
			Positive(),
		field.String("channel_type").
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.JSON("config", map[string]any{}).
			Optional().
			Default(map[string]any{}),
		field.Bool("enable").
			Default(true),
		field.Int("priority").
			Default(0),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (NotificationChannel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("channel_type"),
		index.Fields("user_id", "name").
			Unique(),
	}
}

func (NotificationChannel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("notification_channels").
			Field("user_id").
			Required().
			Unique(),
	}
}
