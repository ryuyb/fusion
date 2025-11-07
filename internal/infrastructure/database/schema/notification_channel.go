package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/ryuyb/fusion/internal/pkg/entgo/mixin"
)

// NotificationChannel holds the schema definition for the NotificationChannel entity.
type NotificationChannel struct {
	ent.Schema
}

// Fields of the NotificationChannel.
func (NotificationChannel) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.Int64("user_id").
			Positive(),
		field.Enum("channel_type").
			Values("email", "webhook", "telegram", "discord", "feishu"),
		field.String("name").
			NotEmpty(),
		field.JSON("config", map[string]interface{}{}).
			Optional(),
		field.Bool("is_enabled").
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

// Indexes of the NotificationChannel.
func (NotificationChannel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "is_enabled"),
		index.Fields("user_id", "priority"),
	}
}

// Edges of the NotificationChannel.
func (NotificationChannel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("notification_channels").
			Field("user_id").
			Required().
			Unique(),
	}
}

// Mixin of the NotificationChannel.
func (NotificationChannel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}
