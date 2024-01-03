package query

type QueryMod interface {
	Apply(q *Query)
}
