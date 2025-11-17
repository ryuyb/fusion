package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// StreamingPlatform stores metadata about supported live-streaming providers.
type StreamingPlatform struct {
	ent.Schema
}

func (StreamingPlatform) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable(),
		field.String("type").
			NotEmpty().
			Unique(),
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("description").
			Optional().
			Nillable(),
		field.String("base_url").
			NotEmpty(),
		field.String("logo_url").
			Optional().
			Nillable(),
		field.Bool("enabled").
			Default(true),
		field.Int("priority").
			Default(0),
		field.JSON("metadata", map[string]string{}).
			Optional().
			Default(map[string]string{}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (StreamingPlatform) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("name"),
	}
}
