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

func Paginate[T any](pt PaginationType, q *Query, d []*T) ([]*T, *PaginationResponse) {
	if pt == OffsetPaginationType {
		return PaginateOffsetPagination(q.OffsetPagination, d)
	}

	return PaginateCursorPagination(q.CursorPagination, d)
}
