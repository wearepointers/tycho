package query

// type PaginationResponse interface {
// 	any
// }

// type Pagination interface {
// 	// Paginate(d any) ([]any, PaginationResponse)

// 	Apply(q *Query)
// 	Mods() []qm.QueryMod
// 	SQL() string
// }

type PaginationResponse struct {
	HasNextPage    bool   `json:"has_next_page"`
	NextPageCursor string `json:"next_page_cursor,omitempty"`
	HasPrevPage    bool   `json:"has_prev_page"`
	PrevPageCursor string `json:"prev_page_cursor,omitempty"`
}
