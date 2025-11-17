package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.String("username").
			NotEmpty(),
		field.String("email").
			NotEmpty(),
		field.String("password").
			NotEmpty().
			Sensitive(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
		index.Fields("email"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("followed_streamers", UserFollowedStreamer.Type),
		edge.To("notification_channels", NotificationChannel.Type),
	}
}
