package query

import (
	"errors"

	"github.com/expanse-agency/tycho/utils"
)

type Order string

var (
	ASC  Order = "asc"
	DESC Order = "desc"
)

var ordering = map[Order]bool{
	ASC:  true,
	DESC: true,
}

func (o Order) IsValid() bool {
	return !ordering[o]
}

// { "name": "asc", "otherfield": "desc"}
type SortMap map[string]Order

type Sort struct {
	Fields []*SortField
}

type SortField struct {
	Field string
	Order Order
}

func ParseSort(raw string) (*Sort, error) {
	sortMap, err := utils.Unmarshal[SortMap](raw)
	if err != nil {
		return nil, err
	}

	if sortMap == nil {
		return nil, errors.New("sort map is nil")
	}

	return parseSortMap(sortMap), nil
}

func parseSortMap(sortMap *SortMap) *Sort {
	var fields []*SortField
	for key, order := range *sortMap {
		if !order.IsValid() {
			continue
		}

		fields = append(fields, &SortField{
			Field: key,
			Order: order,
		})
	}

	return &Sort{
		Fields: fields,
	}
}
