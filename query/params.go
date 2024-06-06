package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
)

////////////////////////////////////////////////////////////////////
// Param
////////////////////////////////////////////////////////////////////

type param struct {
	Column string
	Value  string
}

type ParamSlice []param

func NewParam(column, value string) param {
	return param{Column: column, Value: value}
}

func (p *param) sql(tn string) (string, string) {
	return sql.Where(sql.Column(tn, p.Column), sql.Equal, "?"), p.Value
}

func (p *param) mods(tn string) qm.QueryMod {
	return qm.Where(sql.Where(sql.Column(tn, p.Column), sql.Equal, "?"), p.Value)
}

////////////////////////////////////////////////////////////////////
// Params
////////////////////////////////////////////////////////////////////

type params struct {
	Params []*param
}

func (p *params) Apply(q *Query) {
	q.setParams(p)
}

func (p *params) isEmpty() bool {
	return p == nil || len(p.Params) <= 0
}

func (d *Dialect) ParseParams(raw ...param) *params {
	pms := make([]*param, len(raw))
	for i, param := range raw {
		p := param
		pms[i] = &p
	}

	return &params{pms}
}

func (p *params) SQL(tn string) (string, []any) {
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

func (p *params) Mods(tn string) []qm.QueryMod {
	var and []qm.QueryMod

	for _, param := range p.Params {
		and = append(and, param.mods(tn))
	}

	if len(and) > 1 {
		return []qm.QueryMod{qm.Expr(and...)}
	}

	return and
}
