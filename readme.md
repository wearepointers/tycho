# Tycho - Query filtering, sorting, and pagination for Go APIs

Tycho is a library for filtering, sorting, and paginating queries in Go APIs. You can use it standalone with our own SQL builder or use the query mods for [sqlboiler](https://github.com/volatiletech/sqlboiler).



## Installation

```bash
go get github.com/expanse-agency/tycho
```

## Usage

```go
package main

import (
    "fmt"

   "github.com/expanse-agency/tycho"
   "github.com/gin-gonic/gin"
)


func (s *Service) get(c *gin.Context) {
	selectQuery := query.ParseSelectQuery(c, 100, "id", "name", "url")
	tychoSQL, tychoArgs := selectQuery.SQL("link") // Get the SQL and args via Tycho

	sqlBoilerMods := append(selectQuery.Mods("link"), qm.From("link"))
	sqlBoilerSQL, sqlBoilerArgs := queries.BuildQuery(dm.NewQuery(sqlBoilerMods...)) // Get the SQL and args via SQLboiler

	links, _ := dm.Links(selectQuery.Mods(dm.TableNames.Link)...).All(c, s.db) // Get the links via SQLboiler

	server.Return(c, gin.H{
		"tychoSQL":      tychoSQL,
		"tychoArgs":     tychoArgs,
		"sqlBoilerSQL":  sqlBoilerSQL,
		"sqlBoilerArgs": sqlBoilerArgs,
		"links":         links,
	})
}

// For select queries
func ParseSelectQuery(c *gin.Context, maxLimit int, searchColumns ...string) *query.Query {
	filter := query.ParseFilter(c.Query("filter"))
	sort := query.ParseSort(c.Query("sort"))
	relation := query.ParseRelation(c.Query("relation"))
	search := query.ParseSearch(c.Query("search"), searchColumns)

	pagination := query.ParsePagination(c.Query("pagination"), maxLimit)
	return query.NewQuery(query.Postgres, filter, sort, pagination, relation, search)
}

// For update queries
func ParseUpdateQuery(c *gin.Context) *query.Query {
	filter := query.ParseFilter(c.Query("filter"))
	relation := query.ParseRelation(c.Query("relation"))

	return query.NewQuery(query.Postgres, filter, relation)
}

// For single queries
func ParseSingleQuery(c *gin.Context, column string) *query.Query {
	param := query.ParseParam(column, c.Param(column))
	filter := query.ParseFilter(c.Query("filter"))
	relation := query.ParseRelation(c.Query("relation"))

	return query.NewQuery(query.Postgres, filter, relation, param)
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
https://domain.com/endpoint?sorting={"column": "asc", "column2": "desc"}
```

## Pagination

```
https://domain.com/endpoint?pagination={"page": 1, "limit": 10}
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

MIT © [Expanse Agency](./LICENSE) 2024