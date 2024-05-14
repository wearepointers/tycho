# Tycho - Query filtering, sorting, and pagination for Go APIs

Tycho is a library for filtering, sorting, and paginating queries in Go APIs. You can use it standalone with our own SQL builder or use the query mods for [sqlboiler](https://github.com/volatiletech/sqlboiler).

## TODO
- [ ] Own time format for cursor parsing values
- [ ] More values for cursor like int, float, bool, etc.
- [ ] Include columns in cursor (col:value)
- [ ] Multiple params
- [ ] Fix backward cursor pagination
- [ ] Update pagination docs

## Installation

```bash
go get github.com/wearepointers/tycho
```

## Usage

```go
package main

import (
    "fmt"

   "github.com/wearepointers/tycho/query"
   "github.com/gin-gonic/gin"
)

// To prevent filtering/sorting on columns that don't exist or shouldn't be filtered/sorted on
var TablesWithColumnsMap = map[string]map[string]bool{
	"table_name": {
		"id":   true,
		"name": true,
		"url":  true,
		"tag":  false,
		"domain": true,
		// ...
	},
}

func (s *Service) get(c *gin.Context) {
	// TablesWithColumnsMap can be nil if you want to allow filtering/sorting without checking
	// Search columns can be none if you don't want to allow searching
	selectQuery := ParseListQuery(c, TablesWithColumnsMap[dm.TableNames.TableName], "id", "name", "url")
	tychoSQL, tychoArgs := selectQuery.SQL(dm.TableNames.TableName) // Get the SQL and args via Tycho

	sqlBoilerMods := append(selectQuery.Mods(dm.TableNames.TableName), qm.From(dm.TableNames.TableName)) // Mods is for list, bareMods is for update, count, etc.
	sqlBoilerSQL, sqlBoilerArgs := queries.BuildQuery(dm.NewQuery(sqlBoilerMods...)) // Get the SQL and args via SQLboiler

	links, _ := dm.Links(selectQuery.Mods(dm.TableNames.TableName)...).All(c, s.db) // Get the links via SQLboiler

	paginatedRecords, pagination := query.Paginate(paginationType, pm.Query, links)

	server.Return(c, gin.H{
		"tychoSQL":      tychoSQL,
		"tychoArgs":     tychoArgs,
		"sqlBoilerSQL":  sqlBoilerSQL,
		"sqlBoilerArgs": sqlBoilerArgs,
		"links":         links,
		"records":       paginatedRecords,
		"pagination":    pagination,
	})
}

const (
	maxLimit = 50
	driver   = query.Postgres
	paginationType = query.CursorPaginationType
)

// For list queries, but used with bareMods for single result queries (like sum, count, etc.)
func ParseListQuery(c *gin.Context, tc query.TableColumns, searchColumns ...string) *query.Query {
	filter := query.ParseFilter(c.Query("filter"), tc)
	sort := query.ParseSort(c.Query("sort"), tc)
	relation := query.ParseRelation(c.Query("expand"))
	search := query.ParseSearch(c.Query("search"), searchColumns)

	pagination := query.ParsePagination(c.Query("pagination"), paginationType, maxLimit, false, sort)
	return query.NewQuery(driver, filter, sort, pagination, relation, search)
}

// For single queries
func ParseSingleQuery(c *gin.Context, column string) *query.Query {
	param := query.ParseParam(column, c.Param(column))
	relation := query.ParseRelation(c.Query("expand"))

	return query.NewQuery(query.Postgres, nil, nil, relation, param)
}


```

## Filtering

```
https://domain.com/endpoint?filter={"column": {"operator": "value", "or": {"operator": "value"}}, "or": {"column": {"operator": "value"}}}

```

### Operators

```
eq (equal): any
neq (not equal): any
gt (greater than): number, date
gte (greater than or equal): number, date
lt (less than): number, date
lte (less than or equal): number, date
in (in): any
nin (not in): any
c (contains): string
nc (not contains): string
sw (starts with): string
ew (ends with): string
null (is null): boolean
```

## Sorting

```
https://domain.com/endpoint?sorting=[{"colunn":"name", "order":"ASC"}]
```

## Offset Pagination

```
https://domain.com/endpoint?pagination={"page": 1, "limit": 10}
```

## Cursor Pagination (backward does not work!)

```
https://domain.com/endpoint?pagination={"cursor": "optional cursor", "limit": 10}
```

## Relation

```
https://domain.com/endpoint?relation=["table", "table2"]
```

## Param

``` 
https://domain.com/endpoint/:param
```

## Search

```
https://domain.com/endpoint?search=term
```

## License

MIT © [Pointers](./LICENSE) 2024
