package sql

import (
	"fmt"
	"strings"
)

func QueryEnd(s string) string {
	return fmt.Sprint(s, ";")
}

func Query(c ...string) string {
	return strings.Join(c, " ")
}

func Clause(op LogicalOperator, mods ...string) string {
	return strings.Join(mods, op.String(true))
}

func Expr(s string) string {
	return fmt.Sprint("(", s, ")")
}

func ConvertQuestionMarks(sql string) string {
	count := strings.Count(sql, "?")

	for i := 1; i <= count; i++ {
		sql = strings.Replace(sql, "?", fmt.Sprintf("$%d", i), 1)
	}

	return sql
}

func Column(table, column string) string {
	return fmt.Sprintf(`"%s"."%s"`, table, column)
}

func Group(s ...string) string {
	return strings.Join(s, ", ")
}

func GeneratePlaceholders(count int) []string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}

	return placeholders
}
