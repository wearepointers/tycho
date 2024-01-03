package query

import (
	"github.com/expanse-agency/tycho/sql"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Query struct {
	dialect    *Dialect
	Filter     *Filter
	Sort       *Sort
	Pagination *Pagination
	Relation   *Relation
	Param      *Param
	Search     *Search
}

func NewQuery(d Driver, mods ...QueryMod) *Query {
	q := &Query{}
	q.setDialect(d.Dialect())
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

func (q *Query) setFilter(f *Filter) {
	q.Filter = f
}

func (q *Query) setPagination(p *Pagination) {
	q.Pagination = p
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

	if q.Sort != nil {
		s = append(s, "ORDER BY")
		ss := q.Sort.SQL(tn)

		s = append(s, ss)
	}

	if q.Pagination != nil {
		s = append(s, q.Pagination.SQL())
	}

	if len(s) <= 0 {
		return "", nil
	}

	s = append(s, ";")
	return sql.Query(s...), args

}

func (q *Query) Mods(tn string) []qm.QueryMod {
	var mods []qm.QueryMod

	if q.Filter != nil {
		mods = append(mods, q.Filter.Mods(tn)...)
	}

	if q.Sort != nil {
		mods = append(mods, q.Sort.Mods(tn)...)
	}

	if q.Pagination != nil {
		mods = append(mods, q.Pagination.Mods()...)
	}

	if q.Relation != nil {
		mods = append(mods, q.Relation.Mods()...)
	}

	if q.Param != nil {
		mods = append(mods, q.Param.Mods(tn))
	}

	if q.Search != nil {
		mods = append(mods, q.Search.Mods(tn))
	}

	return mods
}
