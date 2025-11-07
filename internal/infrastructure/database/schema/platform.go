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

// Platform holds the schema definition for the Platform entity.
type Platform struct {
	ent.Schema
}

// Fields of the Platform.
func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty(),
		field.Enum("platform_type").
			Values("douyu", "huya", "bilibili").
			Immutable(),
		field.JSON("config", map[string]interface{}{}).
			Optional(),
		field.Enum("status").
			Values("active", "inactive").
			Default("active"),
		field.Int("poll_interval").
			Default(60).
			Positive(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the Platform.
func (Platform) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform_type").
			Unique(),
	}
}

// Edges of the Platform.
func (Platform) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("streamers", Streamer.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}

// Mixin of the Platform.
func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.SoftDeleteMixin{},
	}
}
