package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/sql"
	"github.com/wearepointers/tycho/utils"
)

////////////////////////////////////////////////////////////////////
// Filter
////////////////////////////////////////////////////////////////////

type filter struct {
	Columns []*filterColumn
	Or      *filter
}

type filterColumn struct {
	Column string
	Where  []*filterColumnWhere
	Or     *filterColumn
}

type filterColumnWhere struct {
	Operator operator
	Value    any
}

func (f *filter) Apply(q *Query) {
	if f == nil {
		return
	}

	q.setFilter(f)
}

func (f *filter) isEmpty() bool {
	return f == nil || len(f.Columns) <= 0 && f.Or == nil
}

type filterMap map[string]json.RawMessage
type filterMapColumn map[operator]json.RawMessage

type ValidatorFunc func(k string) bool

func (d *Dialect) ParseFilter(raw string, validateFunc ValidatorFunc) *filter {
	filterMap, err := utils.Unmarshal[filterMap](raw)
	if err != nil {
		return nil
	}

	return filterMap.parse(validateFunc, d.DBCasing)
}

func (fm *filterMap) parse(validateFunc ValidatorFunc, dbCasing casing) *filter {
	if fm == nil {
		return nil
	}

	var columns []*filterColumn
	var or *filter
	for key, value := range *fm {
		s := string(value)

		key := dbCasing.string(key) // makes key case agnostic

		if operator(key).isOr() {
			filterMap, err := utils.Unmarshal[filterMap](s)
			if err != nil {
				continue
			}

			or = filterMap.parse(validateFunc, dbCasing)
			continue
		}

		// Ability to validate column, for example check if column exists in table
		if validateFunc != nil && !validateFunc(key) {
			continue
		}

		filterMapColumn, err := utils.Unmarshal[filterMapColumn](s)
		if err != nil {
			continue
		}

		columns = append(columns, filterMapColumn.parse(key))
	}

	return &filter{
		Columns: columns,
		Or:      or,
	}
}

func (fmc *filterMapColumn) parse(column string) *filterColumn {
	if fmc == nil {
		return nil
	}

	var where []*filterColumnWhere
	var or *filterColumn

	for operator, value := range *fmc {
		s := string(value)

		if operator.isOr() {
			filterMapColumn, err := utils.Unmarshal[filterMapColumn](s)
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

		if !operator.IsValid(anyValuePointer) {
			continue
		}

		where = append(where, &filterColumnWhere{
			Operator: operator,
			Value:    *anyValuePointer,
		})
	}

	return &filterColumn{
		Column: column,
		Where:  where,
		Or:     or,
	}
}

func (f *filter) SQL(tn string) (string, []any) {
	if f.isEmpty() {
		return "", nil
	}

	var and []string
	var args []any

	for _, c := range f.Columns {
		andSQL, andArgs := c.sql(tn)
		if andSQL == "" {
			continue
		}

		and = append(and, andSQL)
		args = append(args, andArgs...)
	}

	andSQL := sql.Clause(sql.AND, and...)

	if f.Or != nil {
		orSQL, orArgs := f.Or.SQL(tn)
		if orSQL == "" {
			return andSQL, args
		}

		if len(and) > 1 {
			andSQL = sql.Expr(andSQL)
		}

		if len(orArgs) > 1 && (!strings.HasPrefix(orSQL, "(") || !strings.HasSuffix(orSQL, ")")) {
			orSQL = sql.Expr(orSQL)
		}

		return sql.Clause(sql.OR, andSQL, orSQL), append(args, orArgs...)
	}

	return andSQL, args
}

func (f *filter) Mods(tn string) []qm.QueryMod {
	if f.isEmpty() {
		return nil
	}

	var and []qm.QueryMod

	for _, c := range f.Columns {
		and = append(and, c.mods(tn)...)
	}

	if f.Or != nil {
		orMods := f.Or.Mods(tn)
		if orMods == nil {
			return and
		}

		if len(and) > 1 {
			and = []qm.QueryMod{qm.Expr(and...)}
		}

		return append(and, qm.Or2(qm.Expr(orMods...)))

	}

	return and
}

func (f *filterColumn) sql(tn string) (string, []any) {
	var s []string
	var args []any

	for _, w := range f.Where {
		column := sql.Column(tn, f.Column)
		s1, args1 := w.Operator.sql(column, w.Value)
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

func (f *filterColumn) mods(tn string) []qm.QueryMod {
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

		if len(orMods) > 1 {
			return []qm.QueryMod{qm.Expr(append(andMods, orMods...)...)}
		}

		return []qm.QueryMod{qm.Expr(append(andMods, qm.Or2(orMods[0]))...)}
	}

	return andMods
}

////////////////////////////////////////////////////////////////////
// Operator
////////////////////////////////////////////////////////////////////

type operator string

var (
	equal              operator = "eq"
	notEqual           operator = "neq"
	greaterThan        operator = "gt"
	greaterThanOrEqual operator = "gte"
	lessThan           operator = "lt"
	lessThanOrEqual    operator = "lte"
	in                 operator = "in"
	notIn              operator = "nin"
	contains           operator = "c"
	notContains        operator = "nc"
	startsWith         operator = "sw"
	endsWith           operator = "ew"
	null               operator = "null"
	or                 operator = "or"
)

var operators = map[operator]bool{
	equal:              true,
	notEqual:           true,
	greaterThan:        true,
	greaterThanOrEqual: true,
	lessThan:           true,
	lessThanOrEqual:    true,
	in:                 true,
	notIn:              true,
	contains:           true,
	notContains:        true,
	startsWith:         true,
	endsWith:           true,
	null:               true,
	or:                 true,
}

func (o operator) IsValid(v *any) bool {
	if !operators[o] {
		return false
	}

	if v != nil {
		return o.acceptValueKind(v)
	}

	return false
}

func (o operator) acceptValueKind(v *any) bool {
	pV := *v
	switch o {
	case equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual, contains, notContains, startsWith, endsWith:
		if is, _ := utils.IsString(pV); is {
			return true
		}

		if is, _ := utils.IsInt(pV); is {
			return true
		}

		if is, _ := utils.IsInt8(pV); is {
			return true
		}

		if is, _ := utils.IsInt16(pV); is {
			return true
		}

		if is, _ := utils.IsInt32(pV); is {
			return true
		}

		if is, _ := utils.IsInt64(pV); is {
			return true
		}

		if is, _ := utils.IsFloat32(pV); is {
			return true
		}

		if is, _ := utils.IsFloat64(pV); is {
			return true
		}

		return false

	case in, notIn:
		if is, _ := utils.IsSlice[any](pV); is {
			return true
		}

		return false
	case null:
		if is, _ := utils.IsBool(pV); is {
			return true
		}

		return false
	case or:
		return reflect.Map == reflect.TypeOf(v).Kind()
	default:
		return false
	}
}

func (o operator) isOr() bool {
	return o == or
}

func (o operator) sql(c string, v any) (string, any) {
	switch o {
	case equal:
		return sql.Where(c, sql.Equal, "?"), v
	case notEqual:
		return sql.Where(c, sql.NotEqual, "?"), v
	case greaterThan:
		return sql.Where(c, sql.GreaterThan, "?"), v
	case greaterThanOrEqual:
		return sql.Where(c, sql.GreaterThanOrEqual, "?"), v
	case lessThan:
		return sql.Where(c, sql.LessThan, "?"), v
	case lessThanOrEqual:
		return sql.Where(c, sql.LessThanOrEqual, "?"), v
	case in:
		return sql.WhereIn(c, v)
	case notIn:
		return sql.WhereNotIn(c, v)
	case contains:
		return sql.WhereLike(c, v)
	case notContains:
		return sql.WhereNotLike(c, v)
	case startsWith:
		return sql.WhereStartsWith(c, v)
	case endsWith:
		return sql.WhereEndsWith(c, v)
	case null:
		b := v.(bool)
		if b {
			return sql.Where(c, sql.IsNull), nil
		}
		return sql.Where(c, sql.IsNotNull), nil
	default:
		return "", nil
	}
}

func (o operator) mod(c string, v any) qm.QueryMod {
	switch o {
	case equal:
		return qm.Where(sql.Where(c, sql.Equal, "?"), v)
	case notEqual:
		return qm.Where(sql.Where(c, sql.NotEqual, "?"), v)
	case greaterThan:
		return qm.Where(sql.Where(c, sql.GreaterThan, "?"), v)
	case greaterThanOrEqual:
		return qm.Where(sql.Where(c, sql.GreaterThanOrEqual, "?"), v)
	case lessThan:
		return qm.Where(sql.Where(c, sql.LessThan, "?"), v)
	case lessThanOrEqual:
		return qm.Where(sql.Where(c, sql.LessThanOrEqual, "?"), v)
	case in:
		isSlice, val := utils.IsSlice[any](v)
		if !isSlice {
			return nil
		}

		return qm.WhereIn(sql.Query(c, sql.In.String(), "?"), val...)
	case notIn:
		isSlice, val := utils.IsSlice[any](v)
		if !isSlice {
			return nil
		}

		return qm.WhereIn(sql.Query(c, sql.NotIn.String(), "?"), val...)
	case contains:
		return qm.Where(sql.Where(c, sql.Like, "?"), fmt.Sprint("%", v, "%"))
	case notContains:
		return qm.Where(sql.Where(c, sql.NotLike, "?"), fmt.Sprint("%", v, "%"))
	case startsWith:
		return qm.Where(sql.Where(c, sql.Like, "?"), fmt.Sprint(v, "%"))
	case endsWith:
		return qm.Where(sql.Where(c, sql.Like, "?"), fmt.Sprint("%", v))
	case null:
		isBool, val := utils.IsBool(v)
		if !isBool {
			return nil
		}

		if val {
			return qm.Where(sql.Where(c, sql.IsNull))
		}
		return qm.Where(sql.Where(c, sql.IsNotNull))
	default:
		return nil
	}
}
