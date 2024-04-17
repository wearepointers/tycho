package query

import (
	"strings"

	"github.com/wearepointers/tycho/sql"
	"github.com/wearepointers/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// TODO: include sort columns also in cursor, then they don't have to add sorting query anymore but can just use the next cursor to paginate in their sorted list?

// CREATE UNIQUE INDEX ON people (firstname, id) INCLUDE (lastname);

const (
	splitter       = "::"
	columnSplitter = ",,"
)

type Direction string

const (
	forward  Direction = "forward"
	backward Direction = "backward"
)

func (direction Direction) logicalOperator(orderBy sql.Order) sql.Operator {
	if direction == forward {
		if orderBy == sql.ASC {
			return sql.GreaterThan
		}
		return sql.LessThan
	}

	if orderBy == sql.ASC {
		return sql.LessThan
	}

	return sql.GreaterThan
}

///////////////////////////////////////////////////////////////////////
// Cursor
///////////////////////////////////////////////////////////////////////

type Cursor struct {
	columnValues []any // Should be string but any is better for query
	direction    Direction
}

func (c *Cursor) IsBackward() bool {
	if c == nil {
		return false
	}
	return c.direction == backward
}

func (c *Cursor) IsForward() bool {
	if c == nil {
		return false
	}
	return c.direction == forward
}

func (c *Cursor) IsEmpty() bool {
	return c == nil
}

func createCursorFromData[T any](data *T, sortColumns []SortColumn, direction Direction) *Cursor {
	dataMap := utils.NewReflectObject(data)

	columnValues := make([]any, len(sortColumns))
	for i, f := range sortColumns {
		v, err := dataMap.GetUnsafeFieldByTagAsString(f.Column, "boil")
		if err != nil {
			continue
		}
		columnValues[i] = v
	}

	return &Cursor{
		columnValues: columnValues,
		direction:    direction,
	}
}

func (c *Cursor) encode() string {
	if c == nil {
		return ""
	}

	var builder strings.Builder
	for i, v := range c.columnValues {
		val, ok := v.(string)
		if !ok {
			continue
		}
		builder.WriteString(utils.Base64Encode(val))
		if i < len(c.columnValues)-1 {
			builder.WriteString(columnSplitter)
		}
	}

	builder.WriteString(splitter)
	builder.WriteString(string(c.direction))

	return utils.Base64Encode(builder.String())
}

func decodeCursor(s string) *Cursor {
	decoded, err := utils.Base64Decode(s)
	if err != nil {
		return nil
	}

	parts := strings.Split(decoded, splitter)
	if len(parts) != 2 {
		return nil
	}

	columnValuesEncoded := strings.Split(parts[0], columnSplitter)
	direction := Direction(parts[1])

	columnValues := make([]any, len(columnValuesEncoded))
	for i, columnValue := range columnValuesEncoded {
		columnValues[i], err = utils.Base64Decode(columnValue)
		if err != nil {
			return nil
		}
	}

	return &Cursor{
		columnValues: columnValues,
		direction:    direction,
	}
}

///////////////////////////////////////////////////////////////////////
// CursorPagination
///////////////////////////////////////////////////////////////////////

type CursorPaginationInput struct {
	Cursor *string
	Limit  int
}

type CursorPagination struct {
	limit              int
	cursor             *Cursor
	sort               *Sort
	hasAutoIncrementID bool // Uses UUID or auto incrementing ID
}

func (p *CursorPagination) Apply(q *Query) {
	q.setCursorPagination(p)
}

func (p *CursorPagination) Limit() int {
	return p.limit + 1
}

func ParseCursorPagination(raw string, sort *Sort, maxLimit int, hasAutoIncrementID bool) *CursorPagination {
	pagination := &CursorPagination{
		sort:               sort,
		limit:              maxLimit,
		hasAutoIncrementID: hasAutoIncrementID,
	}
	pagination.setDefaultSort()

	paginationInput, err := utils.Unmarshal[CursorPaginationInput](raw)
	if err != nil {
		return pagination
	}

	if paginationInput.Limit > 0 || paginationInput.Limit < maxLimit {
		pagination.limit = paginationInput.Limit
	}

	if paginationInput.Cursor != nil {
		pagination.cursor = decodeCursor(*paginationInput.Cursor)

		// NOTE: weak check to see if the cursor belongs to the sort
		if pagination.cursor != nil && len(pagination.cursor.columnValues) != len(sort.columns) {
			pagination.cursor = nil
		}
	}

	return pagination
}

func PaginateCursorPagination[T any](p *CursorPagination, data []*T) ([]*T, *PaginationResponse) {
	dataCopied := make([]*T, len(data))
	copy(dataCopied, data)

	// When moving forward, we fetch 1 more record to check if there's a next page
	// When moving backward, we now there is a next page
	var hasNextPage = len(data) > p.limit || p.cursor.IsBackward()
	var nextPageCursor *Cursor
	if hasNextPage {
		if p.cursor.IsEmpty() || p.cursor.IsForward() {
			dataCopied = dataCopied[:p.limit]
		}

		if len(dataCopied) > 0 {
			last := dataCopied[len(dataCopied)-1]
			nextPageCursor = createCursorFromData(last, p.sort.columns, forward)
		}
	}

	// When moving backward, we fetch 1 more record to check if there's a prev page
	// When moving forward, we now there is a prev page
	hasPrevPage := !p.cursor.IsEmpty() && (len(data) > p.limit || p.cursor.IsForward())
	var prevPageCursor *Cursor
	if hasPrevPage {
		if p.cursor.IsBackward() {
			dataCopied = dataCopied[1:]
		}

		if len(dataCopied) > 0 {
			first := dataCopied[0]
			prevPageCursor = createCursorFromData(first, p.sort.columns, backward)
		}
	}

	return dataCopied, &PaginationResponse{
		HasNextPage:    hasNextPage,
		NextPageCursor: nextPageCursor.encode(),
		HasPrevPage:    hasPrevPage,
		PrevPageCursor: prevPageCursor.encode(),
	}
}

func (p *CursorPagination) setDefaultSort() {
	p.sort.SetDefault(func(o sql.Order) []SortColumn {
		if p.hasAutoIncrementID {
			return []SortColumn{{Column: "id", Order: o}}
		}

		return []SortColumn{{Column: "created_at", Order: o}, {Column: "id", Order: o}}
	})
}

func (p *CursorPagination) SQL() string {
	return ""
}

// DESC = high to low | latests to oldest | new to old
// ASC = low to high | oldest to latests |  old to new
func (p *CursorPagination) Mods(tn string) []qm.QueryMod {
	m := []qm.QueryMod{
		qm.Limit(p.Limit()),
	}

	m = append(m, p.sort.Mods(tn)...)

	if !p.cursor.IsEmpty() {
		m = append(m, qm.Where(sql.WhereComposite(p.cursor.direction.logicalOperator(p.sort.defaultOrderBy), p.sort.columns, func(c SortColumn) string { return c.Column }), p.cursor.columnValues...))
	}

	return m
}
