package query

import (
	"fmt"
	"reflect"

	"github.com/wearepointers/tycho/sql"
	"github.com/wearepointers/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Operator string

var (
	equal              Operator = "eq"
	notEqual           Operator = "neq"
	greaterThan        Operator = "gt"
	greaterThanOrEqual Operator = "gte"
	lessThan           Operator = "lt"
	lessThanOrEqual    Operator = "lte"
	in                 Operator = "in"
	notIn              Operator = "nin"
	contains           Operator = "c"
	notContains        Operator = "nc"
	startsWith         Operator = "sw"
	endsWith           Operator = "ew"
	null               Operator = "null"
	or                 Operator = "or"
)

var operators = map[Operator]bool{
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
	case equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual, contains, notContains, startsWith, endsWith:
		return reflect.String == v || reflect.Int == v || reflect.Int8 == v || reflect.Int16 == v || reflect.Int32 == v || reflect.Int64 == v || reflect.Float32 == v || reflect.Float64 == v
	case in, notIn:
		return reflect.Slice == v
	case null:
		return reflect.Bool == v
	case or:
		return reflect.Map == v
	default:
		return false
	}
}

func (o Operator) IsOr() bool {
	return o == or
}

func (o Operator) SQL(c string, v any) (string, any) {
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
		return sql.WhereNotIn(c, v)
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

func (o Operator) mod(c string, v any) qm.QueryMod {
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
		// TODO: check if this slice thing still works with ints/strings etc
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
