package schema

import (
	"errors"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("auth_id").
			Unique().
			NotEmpty(),
		field.String("email").
			Unique().
			NotEmpty().
			Validate(func(val string) error {
				if !strings.Contains(val, "@") {
					return errors.New("invalid email")
				}
				return nil
			}),
		field.String("username").
			Unique().
			NotEmpty(),
		field.String("picture").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
		edge.To("owned_pods", Pod.Type),
		edge.From("joined_pods", Pod.Type).
			Ref("users"),
	}
}
