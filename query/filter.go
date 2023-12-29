package query

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/expanse-agency/tycho/sql"
	"github.com/expanse-agency/tycho/utils"
)

type Operator string

var (
	Equal              Operator = "eq"
	NotEqual           Operator = "neq"
	GreaterThan        Operator = "gt"
	GreaterThanOrEqual Operator = "gte"
	LessThan           Operator = "lt"
	LessThanOrEqual    Operator = "lte"
	In                 Operator = "in"
	NotIn              Operator = "nin"
	Contains           Operator = "c"
	NotContains        Operator = "nc"
	StartsWith         Operator = "sw"
	EndsWith           Operator = "ew"
	Null               Operator = "null"
	Or                 Operator = "or"
)

var operators = map[Operator]bool{
	Equal:              true,
	NotEqual:           true,
	GreaterThan:        true,
	GreaterThanOrEqual: true,
	LessThan:           true,
	LessThanOrEqual:    true,
	In:                 true,
	NotIn:              true,
	Contains:           true,
	NotContains:        true,
	StartsWith:         true,
	EndsWith:           true,
	Null:               true,
	Or:                 true,
}

func (o Operator) IsValid(v any) bool {
	if !operators[o] {
		return false
	}

	if v != nil {
		return o.AcceptValueKind(reflect.TypeOf(v).Kind())
	}

	return true
}

func (o Operator) AcceptValueKind(v reflect.Kind) bool {
	switch o {
	case Equal, NotEqual, GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual, Contains, NotContains, StartsWith, EndsWith:
		return reflect.String == v || reflect.Int == v || reflect.Int8 == v || reflect.Int16 == v || reflect.Int32 == v || reflect.Int64 == v || reflect.Float32 == v || reflect.Float64 == v
	case In, NotIn:
		return reflect.Slice == v
	case Null:
		return reflect.Bool == v
	case Or:
		return reflect.Map == v
	default:
		return false
	}
}

func (o Operator) IsOr() bool {
	return o == Or
}

func (o Operator) SQL(c string, v any) (string, any) {
	switch o {
	case Equal:
		return sql.Where(c, sql.Equal, "?"), v
	case NotEqual:
		return sql.Where(c, sql.NotEqual, "?"), v
	case GreaterThan:
		return sql.Where(c, sql.GreaterThan, "?"), v
	case GreaterThanOrEqual:
		return sql.Where(c, sql.GreaterThanOrEqual, "?"), v
	case LessThan:
		return sql.Where(c, sql.LessThan, "?"), v
	case LessThanOrEqual:
		return sql.Where(c, sql.LessThanOrEqual, "?"), v
	case In:
		return sql.WhereIn(c, v)
	case NotIn:
		return sql.WhereNotIn(c, v)
	case Contains:
		return sql.WhereLike(c, v)
	case NotContains:
		return sql.WhereNotIn(c, v)
	case StartsWith:
		return sql.WhereStartsWith(c, v)
	case EndsWith:
		return sql.WhereEndsWith(c, v)
	case Null:
		b := v.(bool)
		if b {
			return sql.Where(c, sql.IsNull), nil
		}
		return sql.Where(c, sql.IsNotNull), nil
	default:
		return "", nil
	}
}

// {"name": {"eq": "test", "or": "test3"}, "age": {"gte": 34, "lte": 65}, "status": {"in": ["active", "paused"]}, "or": {"name":{"eq": "test2"}}}
type FilterMap map[string]json.RawMessage
type FilterMapColumn map[Operator]json.RawMessage

type Filter struct {
	columns []*FilterColumn
	or      *Filter
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

func ParseFilter(raw string) (*Filter, error) {
	filterMap, err := utils.Unmarshal[FilterMap](raw)
	if err != nil {
		return nil, err
	}

	if filterMap == nil {
		return nil, errors.New("filter map is nil")
	}

	return parseFilterMap(filterMap), nil
}

func parseFilterMap(filterMap *FilterMap) *Filter {
	var columns []*FilterColumn
	var or *Filter
	for key, value := range *filterMap {
		s := string(value)
		if Operator(key).IsOr() {
			filterMap, err := utils.Unmarshal[FilterMap](s)
			if err != nil {
				continue
			}

			or = parseFilterMap(filterMap)
			continue
		}

		filterMapColumn, err := utils.Unmarshal[FilterMapColumn](s)
		if err != nil {
			continue
		}

		columns = append(columns, parseFilterMapColumn(filterMapColumn, key))
	}

	return &Filter{
		columns: columns,
		or:      or,
	}
}

func parseFilterMapColumn(filterMapColumn *FilterMapColumn, column string) *FilterColumn {
	var where []*FilterColumnWhere
	var or *FilterColumn

	for operator, value := range *filterMapColumn {
		s := string(value)
		if operator.IsOr() {
			filterMapColumn, err := utils.Unmarshal[FilterMapColumn](s)
			if err != nil {
				continue
			}

			or = parseFilterMapColumn(filterMapColumn, column)
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

// {"name": {"eq": "test", "or": "test3"}, "age": {"gte": 34, "lte": 65}, "status": {"in": ["active", "paused"]}, "or": {"name":{"eq": "test2"}}}
// should prdouce:
// (NAME = 'test' OR NAME = 'test3') AND (AGE >= 34 AND AGE <= 65) AND (STATUS IN ('active', 'paused')) AND (NAME = 'test2')
// But with args
// (NAME = $1 OR NAME = $2) AND (AGE >= $3 AND AGE <= $4) AND (STATUS IN ($5, $6)) AND (NAME = $7)

func (f *Filter) SQL() (string, []any) {
	var s []string
	var args []any

	for _, c := range f.columns {
		s1, args1 := c.SQL()
		if s1 == "" {
			continue
		}

		s = append(s, s1)
		args = append(args, args1...)
	}

	andSQL := sql.Clause(sql.And, s...)

	if f.or != nil {
		andSQL = sql.Expr(andSQL)

		orSQL, orArgs := f.or.SQL()
		if orSQL == "" {
			return andSQL, args
		}

		return sql.Clause(sql.Or, andSQL, orSQL), append(args, orArgs...)
	}

	return andSQL, args
}

func (f *FilterColumn) SQL() (string, []any) {
	var s []string
	var args []any

	for _, w := range f.Where {
		s1, args1 := w.Operator.SQL(f.Column, w.Value)
		if s1 != "" {
			s = append(s, s1)

			if args1 != nil {
				if utils.IsSlice(args1) {
					args = append(args, args1.([]any)...)
					continue
				}

				args = append(args, args1)
			}
		}
	}

	andSQL := sql.Clause(sql.And, s...)
	if len(s) > 1 {
		andSQL = sql.Expr(andSQL)
	}

	if f.Or != nil {
		orSQL, orArgs := f.Or.SQL()
		if orSQL == "" {
			return andSQL, args
		}

		return sql.Expr(sql.Clause(sql.Or, andSQL, orSQL)), append(args, orArgs...)
	}

	return andSQL, args
}

func (f *Filter) PSQL() (string, []any) {
	s, args := f.SQL()
	return sql.Args(s, "$"), args
}
