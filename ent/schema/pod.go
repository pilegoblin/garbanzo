package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Pod holds the schema definition for the Pod entity.
type Pod struct {
	ent.Schema
}

// Fields of the Pod.
func (Pod) Fields() []ent.Field {
	return []ent.Field{
		field.String("pod_name").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now),
		field.String("invite_code").
			Unique(),
	}
}

// Edges of the Pod.
func (Pod) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("owned_pods").
			Unique(),
		edge.To("users", User.Type),
		edge.To("beans", Bean.Type),
	}
}
