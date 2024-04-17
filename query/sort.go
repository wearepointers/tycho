package query

import (
	"github.com/wearepointers/tycho/sql"
	"github.com/wearepointers/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Sort struct {
	filterAllowedOnColumns TableColumns
	columnsInMap           map[string]int
	columns                SortColumnSlice
	defaultOrderBy         sql.Order
}

type SortColumnSlice []SortColumn
type SortColumn struct {
	Column string
	Order  sql.Order
}

func (s *Sort) Apply(q *Query) {
	if s == nil {
		return
	}

	q.setSort(s)
}

// [{"colunn":"name", "order":"ASC"}]

func ParseSort(raw string, allowedColumns TableColumns) *Sort {
	sortColumnSlice, _ := utils.Unmarshal[SortColumnSlice](raw)
	return sortColumnSlice.parse(allowedColumns)
}

func (sortColumnSlice *SortColumnSlice) parse(allowedColumns TableColumns) *Sort {
	var defaultOrderBy = sql.ASC

	if sortColumnSlice == nil {
		return &Sort{filterAllowedOnColumns: allowedColumns, defaultOrderBy: defaultOrderBy}
	}

	var columns []SortColumn
	var columnsInMap = make(map[string]int)

	var i int
	for _, sortColumn := range *sortColumnSlice {
		if !sortColumn.Order.IsValid() || !allowedColumns.Has(sortColumn.Column) {
			continue
		}

		// Duplicate catch
		if _, ok := columnsInMap[sortColumn.Column]; ok {
			continue
		}

		columns = append(columns, sortColumn)
		columnsInMap[sortColumn.Column] = i
		i++
	}

	if len(columns) > 0 {
		defaultOrderBy = columns[0].Order
	}

	return &Sort{
		filterAllowedOnColumns: allowedColumns,
		columns:                columns,
		columnsInMap:           columnsInMap,
		defaultOrderBy:         defaultOrderBy,
	}
}

func (s *Sort) isEmpty() bool {
	return s != nil && len(s.columns) <= 0
}

func (s *Sort) SetDefault(f func(o sql.Order) []SortColumn) {
	columns := f(s.defaultOrderBy)

	if s.isEmpty() {
		s.columns = columns
		return
	}

	// This overwrites any sort
	for _, column := range columns {
		index, ok := s.columnsInMap[column.Column]
		if !ok {
			s.columns = append(s.columns, column)
			continue
		}

		s.columns = append(s.columns[:index], s.columns[index+1:]...) // Remove
		s.columns = append(s.columns, column)                         // Add to end
		s.columnsInMap[column.Column] = len(s.columns) - 1            // Update index
	}
}

func (s *Sort) SQL(tn string) string {
	if len(s.columns) <= 0 {
		return ""
	}

	var orders []string
	for _, f := range s.columns {
		orders = append(orders, sql.Column(tn, f.Column), f.Order.String())
	}

	return sql.Group(orders...)
}

func (s *Sort) Mods(tn string) []qm.QueryMod {
	if len(s.columns) <= 0 {
		return nil
	}

	var mods []qm.QueryMod
	for _, f := range s.columns {
		mods = append(mods, qm.OrderBy(sql.Query(sql.Column(tn, f.Column), f.Order.String())))
	}

	return mods
}
