package sql

type LogicalOperator string

const (
	AND      LogicalOperator = "AND"
	OR       LogicalOperator = "OR"
	ORDER_BY LogicalOperator = "ORDER BY"
)

func (lo LogicalOperator) String(space bool) string {
	if space {
		return " " + string(lo) + " "
	}

	return string(lo)
}

type Operator string

const (
	Equal              Operator = "="
	NotEqual           Operator = "!="
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
	LessThan           Operator = "<"
	LessThanOrEqual    Operator = "<="
	In                 Operator = "IN"
	NotIn              Operator = "NOT IN"
	Like               Operator = "LIKE"
	ILIkE              Operator = "ILIKE"
	NotLike            Operator = "NOT LIKE"
	IsNull             Operator = "IS NULL"
	IsNotNull          Operator = "IS NOT NULL"
)

func (o Operator) String() string {
	return string(o)
}

type Order string

var (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

func (o Order) String() string {
	return string(o)
}

var ordering = map[Order]bool{
	ASC:  true,
	DESC: true,
}

func (o Order) IsValid() bool {
	return ordering[o]
}
