package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Token holds the schema definition for the Token entity.
type Token struct {
	ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("accessToken").Unique(),
		field.String("refreshToken").Unique(),
		field.Time("createdAt").Default(time.Now),
	}
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("tokens").
			Unique(),
	}
}

// Indexes of the Token.
func (Token) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("accessToken").Unique(), // accessToken 필드에 대한 고유 인덱스 추가
	}
}
