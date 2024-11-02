package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title"),
		field.Time("date").
			Default(time.Now()),
		field.String("content"),
		field.Strings("images").Optional(),
		field.String("category").Optional(),
		field.Strings("tag").Optional().StructTag(`json:"tags,omitempty"`).SchemaType(map[string]string{
			"sql": "", // 데이터베이스에 저장하지 않도록 설정
		}),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", Tag.Type).
			Ref("posts"),
	}
}
