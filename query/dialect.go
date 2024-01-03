package query

type Dialect struct {
	UseIndexPlaceholders bool
}

type Driver string

var (
	MySQL    Driver = "mysql"
	Postgres Driver = "postgres"
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
