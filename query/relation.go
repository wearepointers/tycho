package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/utils"
)

// //////////////////////////////////////////////////////////////////
// Relation
// //////////////////////////////////////////////////////////////////
type relation []string

func (r *relation) Apply(q *Query) {
	q.setRelation(r)
}

func (r *relation) isEmpty() bool {
	return r == nil || len(*r) <= 0
}

func (d *Dialect) ParseRelation(raw string) *relation {
	relation, err := utils.Unmarshal[relation](raw)
	if err != nil {
		return nil
	}

	return relation
}

func (r *relation) Mods() []qm.QueryMod {
	if r == nil {
		return nil
	}

	var mods []qm.QueryMod

	for _, relation := range *r {
		mods = append(mods, qm.Load(utils.ToPascalCase(relation), qm.Limit(10)))
	}

	return mods
}
