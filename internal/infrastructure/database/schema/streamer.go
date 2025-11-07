package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/ryuyb/fusion/internal/pkg/entgo/mixin"
)

// Streamer holds the schema definition for the Streamer entity.
type Streamer struct {
	ent.Schema
}

// Fields of the Streamer.
func (Streamer) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.Int64("platform_id").
			Positive(),
		field.String("platform_streamer_id").
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.String("avatar").
			Optional(),
		field.String("description").
			Optional(),
		field.String("room_url").
			Optional(),
		field.Time("last_checked_at").
			Optional().
			Nillable(),
		field.Bool("is_live").
			Default(false),
		field.Time("last_live_at").
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

// Indexes of the Streamer.
func (Streamer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform_id", "platform_streamer_id").
			Unique(),
		index.Fields("is_live"),
		index.Fields("last_checked_at"),
	}
}

// Edges of the Streamer.
func (Streamer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("platform", Platform.Type).
			Ref("streamers").
			Field("platform_id").
			Required().
			Unique(),
		edge.To("followings", UserFollowing.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}

// Mixin of the Streamer.
func (Streamer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}
