package query

import (
	"github.com/expanse-agency/tycho/utils"
	"github.com/volatiletech/sqlboiler/strmangle"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// ["relation", "otherrelation"]
type Relation []string

func (r *Relation) Apply(q *Query) {
	q.setRelation(r)
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
		mods = append(mods, qm.Load(strmangle.TitleCase(relation), qm.Limit(10)))
	}

	return mods
}
