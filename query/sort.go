package query

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
	"github.com/wearepointers/tycho/utils"
)

// //////////////////////////////////////////////////////////////////
// Sort
// //////////////////////////////////////////////////////////////////

type sort struct {
	ColumnsInMap   map[string]int
	Columns        sortColumnSlice
	DefaultOrderBy sql.Order
}

type sortColumnSlice []sortColumn
type sortColumn struct {
	Column string
	Order  sql.Order
}

func (s *sort) Apply(q *Query) {
	if s == nil {
		return
	}

	q.setSort(s)
}

func (s *sort) isEmpty() bool {
	return s == nil || len(s.Columns) <= 0
}

func (d *Dialect) ParseSort(raw string, validateFunc validateColumn) *sort {
	sortColumnSlice, _ := utils.Unmarshal[sortColumnSlice](raw)
	return sortColumnSlice.parse(validateFunc, d.DBCasing)
}

func (scs *sortColumnSlice) parse(validateFunc validateColumn, dbCasing casing) *sort {
	var defaultOrderBy = sql.ASC

	if scs == nil {
		return &sort{DefaultOrderBy: defaultOrderBy}
	}

	var columns []sortColumn
	var columnsInMap = make(map[string]int)

	var i int
	for _, sortColumn := range *scs {
		if !sortColumn.Order.IsValid() || (validateFunc != nil && !validateFunc(sortColumn.Column)) {
			continue
		}

		sortColumn.Column = dbCasing.string(sortColumn.Column) // Makes the column case agnostic

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

	return &sort{
		Columns:        columns,
		ColumnsInMap:   columnsInMap,
		DefaultOrderBy: defaultOrderBy,
	}
}

// func (s *Sort) SetDefault(f func(o sql.Order) []SortColumn) {
// 	columns := f(s.DefaultOrderBy)

// 	if len(s.Columns) <= 0 {
// 		s.Columns = columns
// 		return
// 	}

// 	// This overwrites any sort
// 	for _, column := range columns {
// 		index, ok := s.ColumnsInMap[column.Column]
// 		if !ok {
// 			s.Columns = append(s.Columns, column)
// 			continue
// 		}

// 		s.Columns = append(s.Columns[:index], s.Columns[index+1:]...) // Remove
// 		s.Columns = append(s.Columns, column)                         // Add to end
// 		s.ColumnsInMap[column.Column] = len(s.Columns) - 1            // Update index
// 	}
// }

func (s *sort) SQL(tn string) string {
	if s == nil || len(s.Columns) <= 0 {
		return ""
	}

	var orders []string
	for _, f := range s.Columns {
		orders = append(orders, sql.Query(sql.Column(tn, f.Column), f.Order.String()))
	}

	return sql.Query(string(sql.ORDER_BY), sql.Group(orders...))
}

func (s *sort) Mods(tn string) []qm.QueryMod {
	if s == nil || len(s.Columns) <= 0 {
		return nil
	}

	var mods []qm.QueryMod
	for _, f := range s.Columns {
		mods = append(mods, qm.OrderBy(sql.Query(sql.Column(tn, f.Column), f.Order.String())))
	}

	return mods
}
