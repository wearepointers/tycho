package query

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/expanse-agency/tycho/utils/helpers"
)

type Operator string

var (
	Equal              Operator = "eq"
	NotEqual           Operator = "ne"
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

// {"name": {"eq": "test", "or": "test3"}, "age": {"gte": 34, "lte": 65}, "status": {"in": ["active", "paused"]}, "or": {"name":{"eq": "test2"}}}
type FilterMap map[string]json.RawMessage
type FilterMapField map[Operator]any

type Filter struct {
	Fields []*FilterField `json:"fields"`
	Or     *Filter        `json:"or"`
}

type FilterField struct {
	Field string              `json:"field"`
	Value []*FilterFieldValue `json:"value"`
}

type FilterFieldValue struct {
	Operator Operator `json:"operator"`
	Value    any      `json:"value"`
}

func ParseFilter(raw string) (*Filter, error) {
	filterMap, err := helpers.Unmarshal[FilterMap](raw)
	if err != nil {
		return nil, err
	}

	if filterMap == nil {
		return nil, errors.New("filter map is nil")
	}

	return parseFilterMap(filterMap), nil
}

func parseFilterMap(filterMap *FilterMap) *Filter {
	var fields []*FilterField
	var or *Filter
	for key, value := range *filterMap {
		s := string(value)
		if Operator(key).IsOr() {
			// do something
			filterMap, err := helpers.Unmarshal[FilterMap](s)
			if err != nil {
				continue
			}

			or = parseFilterMap(filterMap)
			continue
		}

		field, err := helpers.Unmarshal[FilterMapField](s)
		if err != nil {
			continue
		}

		var fieldValues []*FilterFieldValue
		for operator, value := range *field {
			if !operator.IsValid(value) {
				continue
			}

			fieldValues = append(fieldValues, &FilterFieldValue{
				Operator: operator,
				Value:    value,
			})
		}

		fields = append(fields, &FilterField{
			Field: key,
			Value: fieldValues,
		})
	}

	return &Filter{
		Fields: fields,
		Or:     or,
	}
}
