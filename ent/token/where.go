// Code generated by ent, DO NOT EDIT.

package token

import (
	"blog/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldID, id))
}

// AccessToken applies equality check predicate on the "accessToken" field. It's identical to AccessTokenEQ.
func AccessToken(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldAccessToken, v))
}

// RefreshToken applies equality check predicate on the "refreshToken" field. It's identical to RefreshTokenEQ.
func RefreshToken(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldRefreshToken, v))
}

// CreatedAt applies equality check predicate on the "createdAt" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCreatedAt, v))
}

// AccessTokenEQ applies the EQ predicate on the "accessToken" field.
func AccessTokenEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldAccessToken, v))
}

// AccessTokenNEQ applies the NEQ predicate on the "accessToken" field.
func AccessTokenNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldAccessToken, v))
}

// AccessTokenIn applies the In predicate on the "accessToken" field.
func AccessTokenIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldAccessToken, vs...))
}

// AccessTokenNotIn applies the NotIn predicate on the "accessToken" field.
func AccessTokenNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldAccessToken, vs...))
}

// AccessTokenGT applies the GT predicate on the "accessToken" field.
func AccessTokenGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldAccessToken, v))
}

// AccessTokenGTE applies the GTE predicate on the "accessToken" field.
func AccessTokenGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldAccessToken, v))
}

// AccessTokenLT applies the LT predicate on the "accessToken" field.
func AccessTokenLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldAccessToken, v))
}

// AccessTokenLTE applies the LTE predicate on the "accessToken" field.
func AccessTokenLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldAccessToken, v))
}

// AccessTokenContains applies the Contains predicate on the "accessToken" field.
func AccessTokenContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldAccessToken, v))
}

// AccessTokenHasPrefix applies the HasPrefix predicate on the "accessToken" field.
func AccessTokenHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldAccessToken, v))
}

// AccessTokenHasSuffix applies the HasSuffix predicate on the "accessToken" field.
func AccessTokenHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldAccessToken, v))
}

// AccessTokenEqualFold applies the EqualFold predicate on the "accessToken" field.
func AccessTokenEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldAccessToken, v))
}

// AccessTokenContainsFold applies the ContainsFold predicate on the "accessToken" field.
func AccessTokenContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldAccessToken, v))
}

// RefreshTokenEQ applies the EQ predicate on the "refreshToken" field.
func RefreshTokenEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldRefreshToken, v))
}

// RefreshTokenNEQ applies the NEQ predicate on the "refreshToken" field.
func RefreshTokenNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldRefreshToken, v))
}

// RefreshTokenIn applies the In predicate on the "refreshToken" field.
func RefreshTokenIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldRefreshToken, vs...))
}

// RefreshTokenNotIn applies the NotIn predicate on the "refreshToken" field.
func RefreshTokenNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldRefreshToken, vs...))
}

// RefreshTokenGT applies the GT predicate on the "refreshToken" field.
func RefreshTokenGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldRefreshToken, v))
}

// RefreshTokenGTE applies the GTE predicate on the "refreshToken" field.
func RefreshTokenGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldRefreshToken, v))
}

// RefreshTokenLT applies the LT predicate on the "refreshToken" field.
func RefreshTokenLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldRefreshToken, v))
}

// RefreshTokenLTE applies the LTE predicate on the "refreshToken" field.
func RefreshTokenLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldRefreshToken, v))
}

// RefreshTokenContains applies the Contains predicate on the "refreshToken" field.
func RefreshTokenContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldRefreshToken, v))
}

// RefreshTokenHasPrefix applies the HasPrefix predicate on the "refreshToken" field.
func RefreshTokenHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldRefreshToken, v))
}

// RefreshTokenHasSuffix applies the HasSuffix predicate on the "refreshToken" field.
func RefreshTokenHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldRefreshToken, v))
}

// RefreshTokenEqualFold applies the EqualFold predicate on the "refreshToken" field.
func RefreshTokenEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldRefreshToken, v))
}

// RefreshTokenContainsFold applies the ContainsFold predicate on the "refreshToken" field.
func RefreshTokenContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldRefreshToken, v))
}

// CreatedAtEQ applies the EQ predicate on the "createdAt" field.
func CreatedAtEQ(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "createdAt" field.
func CreatedAtNEQ(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "createdAt" field.
func CreatedAtIn(vs ...time.Time) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "createdAt" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "createdAt" field.
func CreatedAtGT(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "createdAt" field.
func CreatedAtGTE(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "createdAt" field.
func CreatedAtLT(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "createdAt" field.
func CreatedAtLTE(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldCreatedAt, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Token {
	return predicate.Token(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Token {
	return predicate.Token(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Token) predicate.Token {
	return predicate.Token(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Token) predicate.Token {
	return predicate.Token(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Token) predicate.Token {
	return predicate.Token(sql.NotPredicates(p))
}