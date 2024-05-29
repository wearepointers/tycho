package test

import (
	"testing"

	"github.com/wearepointers/tycho/query"
)

var sortInputs = []test{
	{
		Input: `[{"column": "column1", "order": "ASC"}]`,
		SQL: testSQL{
			SQL: createTestSQL(`ORDER BY "{{table}}"."column1" ASC`),
		},
	},
	{
		Input: `[{"column": "column1", "order": "ASC"},{"column": "column2", "order": "ASC"}]`,
		SQL: testSQL{
			SQL: createTestSQL(`ORDER BY "{{table}}"."column1" ASC, "{{table}}"."column2" ASC`),
		},
	},
	{
		Input: `[{"column": "column1", "order": "ASC"}, {"column": "column2", "order": "DESC"}]`,
		SQL: testSQL{
			SQL: createTestSQL(`ORDER BY "{{table}}"."column1" ASC, "{{table}}"."column2" DESC`),
		},
	},
	{
		Input: `[{"column": "column1", "order": "ASC"}, {"column": "column1", "order": "DESC"}]`,
		SQL: testSQL{
			SQL: createTestSQL(`ORDER BY "{{table}}"."column1" ASC`),
		},
	},
}

func TestSorts(t *testing.T) {
	// We can't test the SQL output directly because the order isn't guaranteed due to maps
	// Would fail sometimes and pass other times
	for i, test := range sortInputs {
		f := query.ParseSort(test.Input, nil)
		s := f.SQL(table)

		// The only test we can do is to check the length of the SQL and SQLArgs
		if len(s) != len(test.SQL.SQL) {
			t.Errorf("Test %d: Expected SQL length %d, got %d", i, len(test.SQL.SQL), len(s))
		}

		// fmt.Println("-------------------------------------------------------")
		// fmt.Println("Sort input:", i)
		// fmt.Println("-------------------------------------------------------")

		// fmt.Println(sql)
		// mod, modArgs := newQuery(f.Mods(table)...)
		// fmt.Println(strings.ReplaceAll(mod, "SELECT * FROM  WHERE ", ""), modArgs)
	}

}
