package query

import "github.com/wearepointers/tycho/utils"

////////////////////////////////////////////////////////////////////
// Driver
////////////////////////////////////////////////////////////////////

type Driver int

const (
	MySQL Driver = iota
	Postgres
)

////////////////////////////////////////////////////////////////////
// Casing
////////////////////////////////////////////////////////////////////

type Casing int

const (
	PascalCase Casing = iota // Represented as PascalCase
	CamelCase                // Represented as camelCase
	SnakeCase                // Represented as snake_case
)

func (c Casing) string(s string) string {
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
	Driver             Driver
	HasAutoIncrementID bool
	APICasing          Casing
	DBCasing           Casing
	PaginationType     PaginationType
	MaxLimit           int
}

func (d *Dialect) Apply(q *Query) {
	if d.Driver == Postgres {
		d.useIndexPlaceholders = true
	}

	q.setDialect(d)
}
