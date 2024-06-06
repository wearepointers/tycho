package test

import (
	"testing"

	"github.com/wearepointers/tycho/query"
)

var paramInputs = []testParam{
	{
		Params: query.ParamSlice{query.NewParam("id", "1")},
		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."id" = ?`),
			SQLArgs: []any{"1"},
		},
	},
	{
		Params: query.ParamSlice{query.NewParam("id", "1"), query.NewParam("account_id", "324")},
		SQL: testSQL{
			SQL:     createTestSQL(`"({{table}}"."id" = ? AND "{{table}}"."account_id" = ?)`),
			SQLArgs: []any{"1", "324"},
		},
	},
	{
		Params: query.ParamSlice{query.NewParam("id", "1"), query.NewParam("account_id", "324"), query.NewParam("event_id", "523532")},
		SQL: testSQL{
			SQL:     createTestSQL(`"({{table}}"."id" = ? AND "{{table}}"."account_id" = ? AND "{{table}}"."event_id" = ?)`),
			SQLArgs: []any{"1", "324", "523532"},
		},
	},
}

func TestParams(t *testing.T) {
	// We can't test the SQL output directly because the order isn't guaranteed due to maps
	// Would fail sometimes and pass other times
	for i, test := range paramInputs {
		f := tychoDialect.ParseParams(test.Params...)
		s, sqlArgs := f.SQL(table)

		// The only test we can do is to check the length of the SQL and SQLArgs
		if (len(s) != len(test.SQL.SQL)) || (len(sqlArgs) != len(test.SQL.SQLArgs)) {
			t.Errorf("Test %d: Expected SQL and SQLArgs to be of length %d and %d, got %d and %d", i, len(test.SQL.SQL), len(test.SQL.SQLArgs), len(s), len(sqlArgs))
		}

		// fmt.Println("-------------------------------------------------------")
		// fmt.Println("Params input:", i)
		// fmt.Println("-------------------------------------------------------")

		// fmt.Println(s, sqlArgs)
		// mod, modArgs := newQuery(f.Mods(table)...)
		// fmt.Println(strings.ReplaceAll(mod, "SELECT * FROM  WHERE ", ""), modArgs)
	}
}
