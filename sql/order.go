package sql

import "strings"

func OrderBy(desc []string, asc []string) string {
	var sql []string
	if len(desc) > 0 {
		sql = append(sql, Query(Group(desc...), DESC.String()))
	}

	if len(asc) > 0 {
		sql = append(sql, Query(Group(asc...), ASC.String()))
	}

	return strings.Join(sql, ",")
}
