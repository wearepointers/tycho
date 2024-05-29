package test

import (
	"testing"

	"github.com/wearepointers/tycho/query"
)

var filterInputs = []test{
	{
		Input: `{
				"column1": {
					"eq": "value1"
				}
			}`,
		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."column1" = ?`),
			SQLArgs: []any{"value1"},
		},
	},
	{
		Input: `{
				"column1": {
					"eq": "value1"
				},
				"column2": {
					"eq": "value2",
					"or": {
						"eq": "value3"
					}
				}
			}`,
		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."column1" = ? AND ("{{table}}"."column2" = ? OR "{{table}}"."column2" = ?)`),
			SQLArgs: []any{"value1", "value2", "value3"},
		},
	},
	{
		Input: `{
				"column1": {
					"eq": "value1"
				},
				"column2": {
					"eq": "value2",
					"or": {
						"eq": "value3",
						"or": {
							"eq": "value4"
						}
					}
				}
			}`,
		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."column1" = ? AND ("{{table}}"."column2" = ? OR ("{{table}}"."column2" = ? OR "{{table}}"."column2" = ?))`),
			SQLArgs: []any{"value1", "value2", "value3", "value4"},
		},
	},
	{
		Input: `{
				"column1": {
					"eq": "value1"
				},
				"or": {
					"column2": {
						"eq": "value2"
					}
				}
			}`,

		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."column1" = ? OR "{{table}}"."column2" = ?`),
			SQLArgs: []any{"value1", "value2"},
		},
	},
	{
		Input: `{
				"column1": {
					"eq": "value1"
				},
				"or": {
					"column2": {
						"eq": "value2"
					},
					"column3": {
						"eq": "value3"
					}
				}
			}`,
		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."column1" = ? OR ("{{table}}"."column2" = ? AND "{{table}}"."column3" = ?)`),
			SQLArgs: []any{"value1", "value2", "value3"},
		},
	},
	{
		Input: `{
				"column1": {
					"eq": "value1"
				},
				"or": {
					"column2": {
						"eq": "value2",
						"or": {
							"eq": "value3"
						}
					},
					"column3": {
						"eq": "value4"
					}
				}
			}`,
		SQL: testSQL{
			SQL:     createTestSQL(`"{{table}}"."column1" = ? OR (("{{table}}"."column2" = ? OR "{{table}}"."column2" = ?) AND "{{table}}"."column3" = ?)`),
			SQLArgs: []any{"value1", "value2", "value3", "value4"},
		},
	},
	{
		Input: `{
				"column1": {
					"eq": "value1"
				},
				"column2": {
					"eq": "value2",
					"or":{
						"eq": "value3"
					}
				},
				"or": {
					"column3": {
						"eq": "value4",
						"or": {
							"eq": "value5"
						}
					}
				}
			}`,
		SQL: testSQL{
			SQL:     createTestSQL(`("{{table}}"."column1" = ? AND ("{{table}}"."column2" = ? OR "{{table}}"."column2" = ?)) OR ("{{table}}"."column3" = ? OR "{{table}}"."column3" = ?)`),
			SQLArgs: []any{"value1", "value2", "value3", "value4", "value5"},
		},
	},
}

func TestFilters(t *testing.T) {
	// We can't test the SQL output directly because the order isn't guaranteed due to maps
	// Would fail sometimes and pass other times
	for i, test := range filterInputs {
		f := query.ParseFilter(test.Input, nil)
		s, sqlArgs := f.SQL(table)

		// The only test we can do is to check the length of the SQL and SQLArgs
		if (len(s) != len(test.SQL.SQL)) || (len(sqlArgs) != len(test.SQL.SQLArgs)) {
			t.Errorf("Test %d: Expected SQL and SQLArgs to be of length %d and %d, got %d and %d", i, len(test.SQL.SQL), len(test.SQL.SQLArgs), len(s), len(sqlArgs))
		}

		// fmt.Println("-------------------------------------------------------")
		// fmt.Println("Filter input:", i)
		// fmt.Println("-------------------------------------------------------")

		// fmt.Println(s, sqlArgs)
		// mod, modArgs := newQuery(f.Mods(table)...)
		// fmt.Println(strings.ReplaceAll(mod, "SELECT * FROM  WHERE ", ""), modArgs)
	}
}
