package query

import "github.com/expanse-agency/tycho/utils"

type Pagination struct {
	Page  int
	Limit int
}

func ParsePagination(raw string, maxLimit int) *Pagination {
	pagination, err := utils.Unmarshal[Pagination](raw)
	if err != nil {
		return &Pagination{
			Page:  0,
			Limit: maxLimit,
		}
	}

	if pagination.Limit > maxLimit {
		pagination.Limit = maxLimit
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 1
	}

	if pagination.Page < 0 {
		pagination.Page = 0
	}

	return pagination
}
