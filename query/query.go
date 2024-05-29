package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
)

type Query struct {
	dialect          *Dialect
	paginationType   PaginationType
	Filter           *Filter
	Params           *Params
	Sort             *Sort
	Relation         *Relation
	OffsetPagination *OffsetPagination
}

type QueryMod interface {
	Apply(q *Query)
}

func NewQuery(dialect Dialect, mods ...QueryMod) *Query {
	q := &Query{}

	dialect.Apply(q)
	for _, mod := range mods {
		mod.Apply(q)
	}

	return q
}

func (q *Query) setDialect(d *Dialect) {
	q.dialect = d
}

func (q *Query) setPaginationType(pt PaginationType) {
	q.paginationType = pt
}

////////////////////////////////////////////////////////////////////
// Set Mods
////////////////////////////////////////////////////////////////////

func (q *Query) setFilter(f *Filter) {
	q.Filter = f
}

func (q *Query) setParams(p *Params) {
	q.Params = p
}

func (q *Query) setSort(s *Sort) {
	q.Sort = s
}

func (q *Query) setRelation(r *Relation) {
	q.Relation = r
}

func (q *Query) setOffsetPagination(p *OffsetPagination) {
	q.OffsetPagination = p
}

////////////////////////////////////////////////////////////////////
// SQL and Mods
////////////////////////////////////////////////////////////////////

func (q *Query) SQL(tn string) (string, []any) {
	var s []string
	var args []any

	if !q.Filter.isEmpty() || !q.Params.isEmpty() {
		s = append(s, "WHERE")
	}

	// Params first because then we can have a params AND filter
	if !q.Params.isEmpty() {
		ps, pa := q.Params.SQL(tn)

		s = append(s, ps)
		args = append(args, pa...)
	}

	if !q.Filter.isEmpty() {
		fs, fa := q.Filter.SQL(tn)

		if len(args) > 0 { // same as if !q.Params.isEmpty() { but little safer
			s = append(s, "AND")
			fs = sql.Expr(fs)
		}

		s = append(s, fs)
		args = append(args, fa...)
	}

	if !q.Sort.isEmpty() {
		ss := q.Sort.SQL(tn)

		s = append(s, ss)
	}

	if len(s) <= 0 {
		return "", nil
	}

	sq := sql.Query(s...)

	if q.dialect.useIndexPlaceholders {
		sq = sql.ConvertQuestionMarks(sq)
	}

	return sql.QueryEnd(sq), args

}

func (q *Query) Mods(tn string) []qm.QueryMod {
	var mods []qm.QueryMod

	if !q.Params.isEmpty() {
		mods = append(mods, q.Params.Mods(tn)...)
	}

	if !q.Filter.isEmpty() {
		if !q.Params.isEmpty() {
			mods = append(mods, qm.Expr(q.Filter.Mods(tn)...))
		}

		if q.Params.isEmpty() {
			mods = append(mods, q.Filter.Mods(tn)...)
		}

	}

	if !q.Sort.isEmpty() {
		mods = append(mods, q.Sort.Mods(tn)...)
	}

	if !q.Relation.isEmpty() {
		mods = append(mods, q.Relation.Mods()...)
	}

	return mods
}
