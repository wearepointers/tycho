package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
)

////////////////////////////////////////////////////////////////////
// Param
////////////////////////////////////////////////////////////////////

type Param struct {
	Column string
	Value  string
}

func (p *Param) sql(tn string) (string, string) {
	return sql.Where(sql.Column(tn, p.Column), sql.Equal, "?"), p.Value
}

func (p *Param) mods(tn string) qm.QueryMod {
	return qm.Where(sql.Where(sql.Column(tn, p.Column), sql.Equal, "?"), p.Value)
}

////////////////////////////////////////////////////////////////////
// Params
////////////////////////////////////////////////////////////////////

type Params struct {
	Params []*Param
}

func (p *Params) Apply(q *Query) {
	q.setParams(p)
}

func (p *Params) isEmpty() bool {
	return p == nil || len(p.Params) <= 0
}

func (d *Dialect) ParseParams(raw ...Param) *Params {
	params := make([]*Param, len(raw))
	for i, param := range raw {
		p := param
		params[i] = &p
	}

	return &Params{params}
}

func (p *Params) SQL(tn string) (string, []any) {
	var and []string
	var args []any

	for _, param := range p.Params {
		paramSQL, paramArgs := param.sql(tn)
		and = append(and, paramSQL)
		args = append(args, paramArgs)
	}

	clause := sql.Clause(sql.AND, and...)

	if len(and) > 1 {
		return sql.Expr(clause), args
	}

	return clause, args
}

func (p *Params) Mods(tn string) []qm.QueryMod {
	var and []qm.QueryMod

	for _, param := range p.Params {
		and = append(and, param.mods(tn))
	}

	if len(and) > 1 {
		return []qm.QueryMod{qm.Expr(and...)}
	}

	return and
}
