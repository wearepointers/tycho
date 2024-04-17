package query

import (
	"encoding/json"

	"github.com/expanse-agency/tycho/sql"
	"github.com/expanse-agency/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Filter struct {
	filterAllowedOnColumns TableColumns
	columns                []*FilterColumn
	or                     *Filter
}

type FilterColumn struct {
	Column string
	Where  []*FilterColumnWhere
	Or     *FilterColumn
}

type FilterColumnWhere struct {
	Operator Operator
	Value    any
}

func (f *Filter) Apply(q *Query) {
	if f == nil {
		return
	}

	q.setFilter(f)
}

// {"name": {"eq": "test", "or": "test3"}, "age": {"gte": 34, "lte": 65}, "status": {"in": ["active", "paused"]}, "or": {"name":{"eq": "test2"}}}
type FilterMap map[string]json.RawMessage
type FilterMapColumn map[Operator]json.RawMessage

func ParseFilter(raw string, allowedColumns TableColumns) *Filter {
	filterMap, err := utils.Unmarshal[FilterMap](raw)
	if err != nil {
		return nil
	}

	return filterMap.parse(allowedColumns)
}

func (filterMap *FilterMap) parse(allowedColumns TableColumns) *Filter {
	if filterMap == nil {
		return nil
	}

	var columns []*FilterColumn
	var or *Filter
	for key, value := range *filterMap {
		if !allowedColumns.Has(key) && !Operator(key).IsOr() {
			continue
		}

		s := string(value)
		if Operator(key).IsOr() {
			filterMap, err := utils.Unmarshal[FilterMap](s)
			if err != nil {
				continue
			}

			or = filterMap.parse(allowedColumns)
			continue
		}

		filterMapColumn, err := utils.Unmarshal[FilterMapColumn](s)
		if err != nil {
			continue
		}

		columns = append(columns, filterMapColumn.parse(key))
	}

	return &Filter{
		filterAllowedOnColumns: allowedColumns,
		columns:                columns,
		or:                     or,
	}
}

func (filterMapColumn *FilterMapColumn) parse(column string) *FilterColumn {
	if filterMapColumn == nil {
		return nil
	}

	var where []*FilterColumnWhere
	var or *FilterColumn

	for operator, value := range *filterMapColumn {
		s := string(value)
		if operator.IsOr() {
			filterMapColumn, err := utils.Unmarshal[FilterMapColumn](s)
			if err != nil {
				continue
			}

			or = filterMapColumn.parse(column)
			continue
		}

		anyValuePointer, err := utils.Unmarshal[any](s)
		if err != nil {
			continue
		}

		if !operator.IsValid(*anyValuePointer) {
			continue
		}

		where = append(where, &FilterColumnWhere{
			Operator: operator,
			Value:    *anyValuePointer,
		})
	}

	return &FilterColumn{
		Column: column,
		Where:  where,
		Or:     or,
	}
}

func (f *Filter) sql(tn string) (string, []any) {
	var s []string
	var args []any

	for _, c := range f.columns {
		s1, args1 := c.sql(tn)
		if s1 == "" {
			continue
		}

		s = append(s, s1)
		args = append(args, args1...)
	}

	andSQL := sql.Clause(sql.AND, s...)

	if f.or != nil {
		andSQL = sql.Expr(andSQL)

		orSQL, orArgs := f.or.sql(tn)
		if orSQL == "" {
			return andSQL, args
		}

		return sql.Clause(sql.OR, andSQL, orSQL), append(args, orArgs...)
	}

	return andSQL, args
}

func (f *Filter) mods(tn string) []qm.QueryMod {
	var andMods []qm.QueryMod

	for _, c := range f.columns {
		cMods := c.mods(tn)
		if cMods == nil {
			continue
		}

		andMods = append(andMods, cMods...)
	}

	if f.or != nil {
		andMods = []qm.QueryMod{qm.Expr(andMods...)}

		orMods := f.or.mods(tn)
		if orMods == nil {
			return andMods
		}

		return append(andMods, orMods...)
	}

	return andMods
}

func (f *FilterColumn) sql(tn string) (string, []any) {
	var s []string
	var args []any

	for _, w := range f.Where {
		column := sql.Column(tn, f.Column)
		s1, args1 := w.Operator.SQL(column, w.Value)
		if s1 != "" {
			s = append(s, s1)

			if args1 != nil {
				isSlice, val := utils.IsSlice[any](args1)
				if isSlice {
					args = append(args, val...)
					continue
				}

				args = append(args, args1)
			}
		}
	}

	andSQL := sql.Clause(sql.AND, s...)
	if len(s) > 1 {
		andSQL = sql.Expr(andSQL)
	}

	if f.Or != nil {
		orSQL, orArgs := f.Or.sql(tn)
		if orSQL == "" {
			return andSQL, args
		}

		return sql.Expr(sql.Clause(sql.OR, andSQL, orSQL)), append(args, orArgs...)
	}

	return andSQL, args
}

func (f *FilterColumn) mods(tn string) []qm.QueryMod {
	var andMods []qm.QueryMod

	for _, w := range f.Where {
		column := sql.Column(tn, f.Column)
		andMods = append(andMods, w.Operator.mod(column, w.Value))
	}

	if len(andMods) > 1 {
		andMods = []qm.QueryMod{qm.Expr(andMods...)}
	}

	if f.Or != nil {
		orMods := f.Or.mods(tn)
		if orMods == nil {
			return andMods
		}

		return []qm.QueryMod{qm.Expr(append(andMods, orMods...)...)}
	}

	return andMods
}

func (f *Filter) SQL(tn string, indexPlaceholders bool) (string, []any) {
	s, args := f.sql(tn)
	if indexPlaceholders {
		return sql.Args(s, "$"), args
	}

	return s, args
}

func (f *Filter) Mods(tn string) []qm.QueryMod {
	return f.mods(tn)
}
