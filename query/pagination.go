package query

import (
	"fmt"

	"github.com/expanse-agency/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Pagination struct {
	Page  int
	Limit int
}

func (p *Pagination) Offset() int {
	return p.Page * p.Limit
}

func (p *Pagination) Apply(q *Query) {
	q.setPagination(p)
}

func ParsePagination(raw string, maxLimit int) *Pagination {
	pagination, err := utils.Unmarshal[Pagination](raw)
	if err != nil {
		return &Pagination{
			Page:  0,
			Limit: maxLimit,
		}
	}

	if pagination.Limit > maxLimit {
		pagination.Limit = maxLimit
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 1
	}

	if pagination.Page < 0 {
		pagination.Page = 0
	}

	return pagination
}

func (p *Pagination) SQL() string {
	offset := p.Offset()

	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", p.Limit)
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, offset)
}

func (p *Pagination) Mods() []qm.QueryMod {
	return []qm.QueryMod{
		qm.Limit(p.Limit),
		qm.Offset(p.Offset()),
	}
}
