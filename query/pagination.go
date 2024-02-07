package query

import (
	"fmt"

	"github.com/expanse-agency/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Pagination struct {
	Page        int
	Limit       int
	limitOffset bool `json:"-"`
}

func (p *Pagination) offset() int {
	return p.Page * p.Limit
}

func (p *Pagination) LimitWithOffset() int {
	if p.limitOffset {
		return p.Limit + 1
	}

	return p.Limit
}

func (p *Pagination) Apply(q *Query) {
	q.setPagination(p)
}

func ParsePagination(raw string, maxLimit int, limitOffset bool) *Pagination {
	pagination, err := utils.Unmarshal[Pagination](raw)
	if err != nil {
		return &Pagination{
			Page:        0,
			Limit:       maxLimit,
			limitOffset: limitOffset,
		}
	}

	if pagination.Limit > maxLimit || pagination.Limit <= 0 {
		pagination.Limit = maxLimit
	}

	if pagination.Page < 0 {
		pagination.Page = 0
	}

	return pagination
}

func (p *Pagination) SQL() string {
	offset := p.offset()

	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", p.LimitWithOffset())
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", p.LimitWithOffset(), offset)
}

func (p *Pagination) Mods() []qm.QueryMod {
	return []qm.QueryMod{
		qm.Limit(p.LimitWithOffset()),
		qm.Offset(p.offset()),
	}
}
