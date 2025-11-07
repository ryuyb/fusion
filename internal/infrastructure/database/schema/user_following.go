package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/ryuyb/fusion/internal/pkg/entgo/mixin"
)

// UserFollowing holds the schema definition for the UserFollowing entity.
type UserFollowing struct {
	ent.Schema
}

// Fields of the UserFollowing.
func (UserFollowing) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.Int64("user_id").
			Positive(),
		field.Int64("streamer_id").
			Positive(),
		field.Bool("notification_enabled").
			Default(true),
		field.Time("last_notified_at").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the UserFollowing.
func (UserFollowing) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "streamer_id").
			Unique(),
		index.Fields("streamer_id"),
		index.Fields("user_id", "notification_enabled"),
	}
}

// Edges of the UserFollowing.
func (UserFollowing) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("followings").
			Field("user_id").
			Required().
			Unique(),
		edge.From("streamer", Streamer.Type).
			Ref("followings").
			Field("streamer_id").
			Required().
			Unique(),
	}
}

// Mixin of the UserFollowing.
func (UserFollowing) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}
