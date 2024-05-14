package query

type PaginationType int

const (
	OffsetPaginationType PaginationType = iota
	CursorPaginationType
)

type PaginationResponse struct {
	HasNextPage    bool   `json:"has_next_page"`
	NextPageCursor string `json:"next_page_cursor,omitempty"`
	HasPrevPage    bool   `json:"has_prev_page"`
	PrevPageCursor string `json:"prev_page_cursor,omitempty"`
}

func Paginate[T any](q *Query, d []*T) ([]*T, *PaginationResponse) {
	if q.paginationType == OffsetPaginationType {
		return paginateOffsetPagination(q.OffsetPagination, d)
	}

	return paginateCursorPagination(q.CursorPagination, d)
}

func ParsePagination(raw string, pt PaginationType, maxLimit int, hasAutoIncrementID bool, sort *Sort) QueryMod {
	if pt == OffsetPaginationType {
		return parseOffsetPagination(raw, maxLimit)
	}

	return parseCursorPagination(raw, maxLimit, hasAutoIncrementID, sort)
}
