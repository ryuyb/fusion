package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserFollowedStreamer stores follow relationship and preferences between users and streamers.
type UserFollowedStreamer struct {
	ent.Schema
}

func (UserFollowedStreamer) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.Int64("user_id").
			Positive(),
		field.Int64("streamer_id").
			Positive(),
		field.String("alias").
			Optional().
			Nillable(),
		field.String("notes").
			Optional().
			Nillable(),
		field.Bool("notifications_enabled").
			Default(true),
		field.JSON("notification_channel_ids", []int64{}).
			Optional().
			Default([]int64{}),
		field.Time("last_notification_sent_at").
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

func (UserFollowedStreamer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("streamer_id"),
		index.Fields("user_id", "streamer_id").
			Unique(),
	}
}

func (UserFollowedStreamer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("followed_streamers").
			Field("user_id").
			Required().
			Unique(),
		edge.From("streamer", Streamer.Type).
			Ref("followers").
			Field("streamer_id").
			Required().
			Unique(),
	}
}
