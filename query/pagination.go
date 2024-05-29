package query

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/utils"
)

////////////////////////////////////////////////////////////////////
// Pagination
////////////////////////////////////////////////////////////////////

type PaginationType int

const (
	OffsetPaginationType PaginationType = iota
	CursorPaginationType
)

type PaginationResponse = map[string]any

func Paginate[T any](q *Query, data []*T) ([]*T, *PaginationResponse) {
	if q.paginationType == OffsetPaginationType {
		return paginateOffsetPagination(q.OffsetPagination, q.dialect.APICasing, data)
	}

	return nil, nil
}

func (d *Dialect) ParsePagination(raw string) QueryMod {
	if d.PaginationType == OffsetPaginationType {
		return d.parseOffsetPagination(raw)
	}

	return nil
}

// //////////////////////////////////////////////////////////////////
// Offset Pagination
// //////////////////////////////////////////////////////////////////
type OffsetPagination struct {
	Page  int
	Limit int
}

func (p *OffsetPagination) Apply(q *Query) {
	q.setPaginationType(OffsetPaginationType)
	q.setOffsetPagination(p)
}

func (p *OffsetPagination) offset() int {
	return p.Page * p.Limit
}

func (p *OffsetPagination) limit() int {
	return p.Limit + 1
}

func paginateOffsetPagination[T any](p *OffsetPagination, apiC Casing, data []*T) ([]*T, *PaginationResponse) {
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

	// TODO: find another way to do this
	r := utils.OmitemptyMap(PaginationResponse{
		apiC.string("hasNextPage"):    hasNextPage,
		apiC.string("nextPageCursor"): nextPageCursor,
		apiC.string("hasPrevPage"):    hasPrevPage,
		apiC.string("prevPageCursor"): prevPageCursor,
	})

	return cData, &r
}

func (d *Dialect) parseOffsetPagination(raw string) *OffsetPagination {
	pagination, err := utils.Unmarshal[OffsetPagination](raw)
	if err != nil {
		return &OffsetPagination{
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
