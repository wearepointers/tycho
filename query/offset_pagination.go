package query

import (
	"fmt"

	"github.com/wearepointers/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type OffsetPagination struct {
	Page  int
	Limit int
}

func (p *OffsetPagination) Apply(q *Query) {
	q.setOffsetPagination(p)
}

func (p *OffsetPagination) offset() int {
	return p.Page * p.Limit
}

func (p *OffsetPagination) limit() int {
	return p.Limit + 1
}

func PaginateOffsetPagination[T any](p *OffsetPagination, d []T) ([]T, *PaginationResponse) {
	var len = len(d)
	var cData = d

	if len >= p.Limit {
		cData = cData[:p.Limit]
	}

	return cData, &PaginationResponse{
		HasNextPage: len > (p.Limit),
		HasPrevPage: len > 0 && p.Page > 0,
	}
}

func ParseOffsetPagination(raw string, maxLimit int) *OffsetPagination {
	pagination, err := utils.Unmarshal[OffsetPagination](raw)
	if err != nil {
		return &OffsetPagination{
			Page:  0,
			Limit: maxLimit,
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

func (p *OffsetPagination) SQL() string {
	offset := p.offset()

	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", p.limit())
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", p.limit(), offset)
}

func (p *OffsetPagination) Mods() []qm.QueryMod {
	return []qm.QueryMod{
		qm.Limit(p.limit()),
		qm.Offset(p.offset()),
	}
}
