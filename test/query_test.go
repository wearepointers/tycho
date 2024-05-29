package test

import (
	"strings"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/wearepointers/tycho/query"
)

type test struct {
	Input string
	SQL   testSQL
}

type testParam struct {
	Params []query.Param
	SQL    testSQL
}

type testSQL struct {
	SQL     string
	SQLArgs []any
}

var dialect = drivers.Dialect{
	LQ: 0x22,
	RQ: 0x22,

	UseIndexPlaceholders:    false,
	UseLastInsertID:         false,
	UseSchema:               false,
	UseDefaultKeyword:       true,
	UseAutoColumns:          false,
	UseTopClause:            false,
	UseOutputClause:         false,
	UseCaseWhenExistsClause: false,
}

var table = "table_name"

// NewQuery initializes a new Query using the passed in QueryMods
func newQuery(mods ...qm.QueryMod) (string, []interface{}) {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return queries.BuildQuery(q)
}

func createTestSQL(sql string) string {
	return strings.ReplaceAll(sql, "{{table}}", table)
}

func TestQuery(t *testing.T) {
	dialect := query.Dialect{
		Driver:             query.MySQL,
		HasAutoIncrementID: false,
		APICasing:          query.CamelCase,
		DBCasing:           query.SnakeCase,
	}

	filter := query.ParseFilter(filterInputs[1].Input, nil)
	params := query.ParseParams(paramInputs[0].Params...)
	sort := query.ParseSort(sortInputs[0].Input, nil)
	q := query.NewQuery(dialect, filter, params, sort)

	s, args := q.SQL(table)

	expectedSQL := createTestSQL(`WHERE "{{table}}"."id" = ? AND ("{{table}}"."column1" = ? AND ("{{table}}"."column2" = ? OR "{{table}}"."column2" = ?)) ORDER BY "{{table}}"."column1" ASC;`)
	expectedArgs := []interface{}{"1", "value1", "value2", "value3"}

	if len(s) != len(expectedSQL) || len(args) != len(expectedArgs) {
		t.Errorf("Expected SQL and SQLArgs to be of length %d and %d, got %d and %d", len(expectedSQL), len(expectedArgs), len(s), len(args))
	}

	// fmt.Println("-------------------------------------------------------")
	// fmt.Println("Query input)
	// fmt.Println("-------------------------------------------------------")

	// mod, modArgs := newQuery(q.Mods(table)...)
	// fmt.Println(s, args)
	// fmt.Println(strings.ReplaceAll(mod, "SELECT * FROM  ", ""), modArgs)

}
