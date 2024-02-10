package query

type Dialect struct {
	UseIndexPlaceholders bool
}

type Driver int

const (
	MySQL Driver = iota
	Postgres
)

func (d Driver) Dialect() *Dialect {
	switch d {
	case MySQL:
		return &Dialect{
			UseIndexPlaceholders: false,
		}
	case Postgres:
		return &Dialect{
			UseIndexPlaceholders: true,
		}
	default:
		return &Dialect{}
	}
}
