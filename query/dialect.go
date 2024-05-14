package query

type Dialect struct {
	UseIndexPlaceholders bool
	HasAutoIncrementID   bool
}

type Driver int

const (
	MySQL Driver = iota
	Postgres
)

func (d Driver) Dialect(hasAutoIncrementID bool) *Dialect {
	switch d {
	case MySQL:
		return &Dialect{
			UseIndexPlaceholders: false,
			HasAutoIncrementID:   hasAutoIncrementID,
		}
	case Postgres:
		return &Dialect{
			UseIndexPlaceholders: true,
			HasAutoIncrementID:   hasAutoIncrementID,
		}
	default:
		return &Dialect{}
	}
}
