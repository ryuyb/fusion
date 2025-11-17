package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Streamer describes a content creator inside supported streaming platforms.
type Streamer struct {
	ent.Schema
}

func (Streamer) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.String("platform_type").
			NotEmpty(),
		field.String("platform_streamer_id").
			NotEmpty(),
		field.String("display_name").
			NotEmpty(),
		field.String("avatar_url").
			Optional().
			Nillable(),
		field.String("room_url").
			Optional().
			Nillable(),
		field.String("bio").
			Optional().
			Nillable(),
		field.JSON("tags", []string{}).
			Optional().
			Default([]string{}),
		field.Time("last_synced_at").
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

func (Streamer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform_type", "platform_streamer_id").
			Unique(),
		index.Fields("display_name"),
	}
}

func (Streamer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("followers", UserFollowedStreamer.Type),
	}
}
