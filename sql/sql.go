package sql

import (
	"fmt"
	"strings"

	"github.com/wearepointers/tycho/utils"
)

func Query(c ...string) string {
	return strings.Join(c, " ")
}

func Clause(op LogicalOperator, mods ...string) string {
	return strings.Join(mods, op.String(true))
}

func Expr(s string) string {
	return fmt.Sprint("(", s, ")")
}

func Args(sql, prefix string) string {
	count := strings.Count(sql, "?")

	for i := 1; i <= count; i++ {
		sql = strings.Replace(sql, "?", fmt.Sprintf("%v%d", prefix, i), 1)
	}

	return sql
}

func Column(c ...string) string {
	return utils.Join(c, ".", `"`, `"`)
}

func Group(s ...string) string {
	return strings.Join(s, ",")
}
