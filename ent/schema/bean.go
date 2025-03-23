package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Bean holds the schema definition for the Bean entity.
type Bean struct {
	ent.Schema
}

// Fields of the Bean.
func (Bean) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
	}
}

// Edges of the Bean.
func (Bean) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
	}
}
