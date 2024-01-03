package query

import (
	"github.com/expanse-agency/tycho/sql"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Param struct {
	column string
	value  any
}

func (p *Param) Apply(q *Query) {
	q.setParam(p)
}

func ParseParam(column string, value any) *Param {
	return &Param{column, value}
}

func (p *Param) SQL(tn string) (string, any) {
	return sql.Where(sql.Column(tn, p.column), sql.Equal, "?"), p.value
}

func (p *Param) Mods(tn string) qm.QueryMod {
	return qm.Where(sql.Where(sql.Column(tn, p.column), sql.Equal, "?"), p.value)
}
