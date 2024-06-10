package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
)

type Query struct {
	dialect          *Dialect
	paginationType   paginationType
	Filter           *filter
	Params           *params
	Sort             *sort
	Relation         *relation
	OffsetPagination *offsetPagination
}

type queryMod interface {
	Apply(q *Query)
}

func (d *Dialect) NewQuery(mods ...queryMod) *Query {
	q := &Query{}
	d.Apply(q)

	for _, mod := range mods {
		mod.Apply(q)
	}

	return q
}

func (q *Query) setDialect(d *Dialect) {
	q.dialect = d
}

func (q *Query) setPaginationType(pt paginationType) {
	q.paginationType = pt
}

////////////////////////////////////////////////////////////////////
// Set Mods
////////////////////////////////////////////////////////////////////

func (q *Query) setFilter(f *filter) {
	q.Filter = f
}

func (q *Query) setParams(p *params) {
	q.Params = p
}

func (q *Query) setSort(s *sort) {
	q.Sort = s
}

func (q *Query) setRelation(r *relation) {
	q.Relation = r
}

func (q *Query) setOffsetPagination(p *offsetPagination) {
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

	if q.OffsetPagination != nil {
		s = append(s, q.OffsetPagination.SQL())
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

	if q.OffsetPagination != nil {
		mods = append(mods, q.OffsetPagination.Mods()...)
	}

	if !q.Relation.isEmpty() {
		mods = append(mods, q.Relation.Mods()...)
	}

	return mods
}
