package query

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/utils"
)

////////////////////////////////////////////////////////////////////
// Pagination
////////////////////////////////////////////////////////////////////

type paginationType int

const (
	OffsetPagination paginationType = iota
	CursorPagination
)

type PaginationResponse = map[string]any

func Paginate[T any](q *Query, data []*T) ([]*T, *PaginationResponse) {
	if q.paginationType == OffsetPagination {
		return paginateOffsetPagination(q.OffsetPagination, q.dialect.APICasing, data)
	}

	return nil, nil
}

func (d *Dialect) ParsePagination(raw string) queryMod {
	if d.PaginationType == OffsetPagination {
		return d.parseOffsetPagination(raw)
	}

	return nil
}

// //////////////////////////////////////////////////////////////////
// Offset Pagination
// //////////////////////////////////////////////////////////////////
type offsetPagination struct {
	Page  int
	Limit int
}

func (p *offsetPagination) Apply(q *Query) {
	q.setPaginationType(OffsetPagination)
	q.setOffsetPagination(p)
}

func (p *offsetPagination) offset() int {
	return p.Page * p.Limit
}

func (p *offsetPagination) limit() int {
	return p.Limit + 1
}

func paginateOffsetPagination[T any](p *offsetPagination, apiC casing, data []*T) ([]*T, *PaginationResponse) {
	var len = len(data)
	var cData = data

	if len >= p.Limit {
		cData = cData[:p.Limit]
	}

	hasNextPage := len > (p.Limit)
	hasPrevPage := len > 0 && p.Page > 0
	nextPageCursor := fmt.Sprint(p.Page + 1)
	prevPageCursor := fmt.Sprint(p.Page - 1)

	if !hasNextPage {
		nextPageCursor = ""
	}

	if !hasPrevPage {
		prevPageCursor = ""
	}

	return cData, &PaginationResponse{
		apiC.string("hasNextPage"):    hasNextPage,
		apiC.string("nextPageCursor"): nextPageCursor,
		apiC.string("hasPrevPage"):    hasPrevPage,
		apiC.string("prevPageCursor"): prevPageCursor,
	}
}

func (d *Dialect) parseOffsetPagination(raw string) *offsetPagination {
	pagination, err := utils.Unmarshal[offsetPagination](raw)
	if err != nil {
		return &offsetPagination{
			Page:  0,
			Limit: d.MaxLimit,
		}
	}

	if pagination.Limit > d.MaxLimit || pagination.Limit <= 0 {
		pagination.Limit = d.MaxLimit
	}

	if pagination.Page < 0 {
		pagination.Page = 0
	}

	return pagination
}

func (p *offsetPagination) SQL() string {
	offset := p.offset()

	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", p.limit())
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", p.limit(), offset)
}

func (p *offsetPagination) Mods() []qm.QueryMod {
	return []qm.QueryMod{
		qm.Limit(p.limit()),
		qm.Offset(p.offset()),
	}
}
