package sql

import "fmt"

type LogicalOperator string

const (
	And LogicalOperator = "AND"
	Or  LogicalOperator = "OR"
)

func (lo LogicalOperator) String() string {
	return fmt.Sprint(" ", string(lo), " ")
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
	NotLike            Operator = "NOT LIKE"
	IsNull             Operator = "IS NULL"
	IsNotNull          Operator = "IS NOT NULL"
)

func (o Operator) String() string {
	return string(o)
}
