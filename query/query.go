package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
)

type TableColumns map[string]bool

func (c TableColumns) Has(column string) bool {
	return c[column]
}

type Query struct {
	dialect          *Dialect
	paginationType   PaginationType
	Filter           *Filter
	Sort             *Sort
	OffsetPagination *OffsetPagination
	CursorPagination *CursorPagination
	Relation         *Relation
	Param            *Param
	Search           *Search
}

func NewQuery(d Driver, pt PaginationType, hasAutoIncrementID bool, mods ...QueryMod) *Query {
	q := &Query{}
	q.setDialect(d.Dialect(hasAutoIncrementID))
	q.setPaginationType(pt)
	q.applyMods(mods...)

	return q
}

func (q *Query) applyMods(mods ...QueryMod) {
	for _, mod := range mods {
		mod.Apply(q)
	}
}

func (q *Query) setDialect(d *Dialect) {
	q.dialect = d
}

func (q *Query) setPaginationType(pt PaginationType) {
	q.paginationType = pt
}

func (q *Query) setFilter(f *Filter) {
	q.Filter = f
}

func (q *Query) setOffsetPagination(cp *OffsetPagination) {
	q.OffsetPagination = cp
}

func (q *Query) setCursorPagination(cp *CursorPagination) {
	q.CursorPagination = cp
}

func (q *Query) setSort(s *Sort) {
	q.Sort = s
}

func (q *Query) setRelation(r *Relation) {
	q.Relation = r
}

func (q *Query) setParam(p *Param) {
	q.Param = p
}

func (q *Query) setSearch(s *Search) {
	q.Search = s
}

func (q *Query) SQL(tn string) (string, []any) {
	var s []string
	var args []any

	if q.Filter != nil || q.Param != nil || q.Search != nil {
		s = append(s, "WHERE")
	}

	if q.Filter != nil {
		fs, fa := q.Filter.SQL(tn, q.dialect.UseIndexPlaceholders)

		s = append(s, fs)
		args = append(args, fa...)
	}

	if q.Param != nil {
		if len(s) > 0 {
			s = append(s, "AND")
		}
		ps, pa := q.Param.SQL(tn)
		s = append(s, ps)
		args = append(args, pa)
	}

	if q.Search != nil {
		if len(s) > 0 {
			s = append(s, "AND")
		}
		ss, sa := q.Search.SQL(tn)
		s = append(s, ss)
		args = append(args, sa)
	}

	// Disable sorting here, we do that in cursor pagination if set
	if !q.Sort.isEmpty() && q.CursorPagination == nil {
		ss := q.Sort.SQL(tn)

		s = append(s, ss)
	}

	if q.OffsetPagination != nil {
		s = append(s, q.OffsetPagination.SQL())
	}

	if q.CursorPagination != nil && q.OffsetPagination == nil {
		s = append(s, q.CursorPagination.SQL())
	}

	if len(s) <= 0 {
		return "", nil
	}

	s = append(s, ";")
	return sql.Query(s...), args

}

func (q *Query) BareMods(tn string) []qm.QueryMod {
	var mods []qm.QueryMod

	if q.Filter != nil {
		mods = append(mods, q.Filter.Mods(tn)...)
	}

	if q.Relation != nil {
		mods = append(mods, q.Relation.Mods()...)
	}

	if q.Param != nil {
		mods = append(mods, q.Param.Mods(tn))
	}

	return mods
}

func (q *Query) Mods(tn string) []qm.QueryMod {
	mods := q.BareMods(tn)

	// Disable sorting here, we do that in cursor pagination if set
	if !q.Sort.isEmpty() && q.CursorPagination == nil {
		mods = append(mods, q.Sort.Mods(tn)...)
	}

	if q.OffsetPagination != nil {
		mods = append(mods, q.OffsetPagination.Mods()...)
	}

	if q.CursorPagination != nil && q.OffsetPagination == nil {
		mods = append(mods, q.CursorPagination.Mods(tn)...)
	}

	if q.Search != nil {
		mods = append(mods, q.Search.Mods(tn))
	}

	return mods
}
