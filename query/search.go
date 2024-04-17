package query

import (
	"fmt"

	"github.com/wearepointers/tycho/sql"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Search struct {
	value   any
	columns []string
}

func (s *Search) Apply(q *Query) {
	q.setSearch(s)
}

func ParseSearch(value string, columns []string) *Search {
	if len(columns) <= 0 || value == "" {
		return nil
	}

	return &Search{value, columns}
}

func (s *Search) SQL(tn string) (string, any) {
	var q []string
	var args []any

	for _, column := range s.columns {
		q = append(q, sql.Where(sql.Column(tn, column), sql.ILIkE, "?"))
	}

	return sql.Expr(sql.Clause(sql.OR, q...)), args
}

func (s *Search) Mods(tn string) qm.QueryMod {
	var mods []qm.QueryMod
	for i, column := range s.columns {
		if i > 0 {
			mods = append(mods, qm.Or2(qm.Where(sql.Where(sql.Column(tn, column), sql.ILIkE, "?"), fmt.Sprint("%", s.value, "%"))))
			continue
		}

		mods = append(mods, qm.Where(sql.Where(sql.Column(tn, column), sql.ILIkE, "?"), fmt.Sprint("%", s.value, "%")))
	}

	return qm.Expr(mods...)
}
