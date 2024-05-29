package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/utils"
)

// //////////////////////////////////////////////////////////////////
// Relation
// //////////////////////////////////////////////////////////////////
type Relation []string

func (r *Relation) Apply(q *Query) {
	q.setRelation(r)
}

func (r *Relation) isEmpty() bool {
	return r == nil || len(*r) <= 0
}

func ParseRelation(raw string) *Relation {
	relation, err := utils.Unmarshal[Relation](raw)
	if err != nil {
		return nil
	}

	return relation
}

func (r *Relation) Mods() []qm.QueryMod {
	if r == nil {
		return nil
	}

	var mods []qm.QueryMod

	for _, relation := range *r {
		mods = append(mods, qm.Load(utils.ToPascalCase(relation), qm.Limit(10)))
	}

	return mods
}
