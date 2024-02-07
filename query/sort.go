package query

import (
	"github.com/expanse-agency/tycho/sql"
	"github.com/expanse-agency/tycho/utils"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Sort struct {
	tableColumns TableColumns
	fields       []*SortField
}

type SortField struct {
	Field string
	Order Order
}

func (s *Sort) Apply(q *Query) {
	if s == nil {
		return
	}

	s.tableColumns = q.sortingAllowedOnColumns
	q.setSort(s)
}

// { "name": "asc", "otherfield": "desc"}
type SortMap map[string]Order

func ParseSort(raw string) *Sort {
	sortMap, err := utils.Unmarshal[SortMap](raw)
	if err != nil {
		return nil
	}

	return sortMap.parse()
}

func (sortMap *SortMap) parse() *Sort {
	if sortMap == nil {
		return nil
	}

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
		fields: fields,
	}
}

func (s *Sort) SQL(tn string) string {
	if len(s.fields) <= 0 {
		return ""
	}

	var orderDesc []string
	var orderAsc []string

	for _, f := range s.fields {
		if !s.tableColumns.Has(f.Field) {
			continue
		}

		if f.Order == desc {
			orderDesc = append(orderDesc, sql.Column(tn, f.Field))
		}

		if f.Order == asc {
			orderAsc = append(orderAsc, sql.Column(tn, f.Field))
		}
	}

	return sql.OrderBy(orderDesc, orderAsc)
}

func (s *Sort) Mods(tn string) []qm.QueryMod {
	if len(s.fields) <= 0 {
		return nil
	}

	var mods []qm.QueryMod

	for _, f := range s.fields {
		if !s.tableColumns.Has(f.Field) {
			continue
		}

		if f.Order == desc {
			mods = append(mods, qm.OrderBy(sql.Query(sql.Column(tn, f.Field), sql.DESC.String())))
		}

		if f.Order == asc {
			mods = append(mods, qm.OrderBy(sql.Query(sql.Column(tn, f.Field), sql.ASC.String())))
		}
	}

	return mods
}
