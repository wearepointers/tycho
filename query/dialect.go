package query

import "github.com/wearepointers/tycho/utils"

////////////////////////////////////////////////////////////////////
// Driver
////////////////////////////////////////////////////////////////////

type driver int

const (
	MySQL driver = iota
	Postgres
)

////////////////////////////////////////////////////////////////////
// Casing
////////////////////////////////////////////////////////////////////

type casing int

const (
	PascalCase casing = iota // Represented as PascalCase
	CamelCase                // Represented as camelCase
	SnakeCase                // Represented as snake_case
)

func (c casing) string(s string) string {
	switch c {
	case PascalCase:
		return utils.ToPascalCase(s)
	case CamelCase:
		return utils.ToCamelCase(s)
	case SnakeCase:
		return utils.ToSnakeCase(s)
	default:
		return "Unknown"
	}
}

////////////////////////////////////////////////////////////////////
// Dialect
////////////////////////////////////////////////////////////////////

type Dialect struct {
	useIndexPlaceholders bool
	// Exported fields
	Driver             driver
	HasAutoIncrementID bool
	APICasing          casing
	DBCasing           casing
	PaginationType     paginationType
	MaxLimit           int
}

func (d *Dialect) Apply(q *Query) {
	if d.Driver == Postgres {
		d.useIndexPlaceholders = true
	}

	q.setDialect(d)
}
