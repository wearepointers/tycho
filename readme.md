# Tycho - Query filtering, sorting, and pagination for Go APIs

Tycho is a library for filtering, sorting, and paginating queries in Go APIs. You can use it standalone with our own SQL builder or use the query mods for [sqlboiler](https://github.com/volatiletech/sqlboiler).

## TODO
- [x] Multiple params
- [x] Case agnostic (snake, camel,pascal, etc.) 
- [ ] Update docs
- [ ] Implement cursor pagination
  - [ ] Own time format for cursor parsing values | remove constant
  - [ ] More values for cursor like int, float, bool, etc.
  - [ ] Include columns in cursor (col:value)
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

// Place this at 1 place in your code.
var dialect = query.Dialect{
	Driver:             query.Postgres,
	HasAutoIncrementID: false,
	APICasing:          query.CamelCase,
	DBCasing:           query.SnakeCase,
	PaginationType:     query.OffsetPagination,
	MaxLimit:           10,
}

// GET /events
func (r *Router) list(c *gin.Context) {
	filter := dialect.ParseFilter(c.Query("filter"), nil)
	sort := dialect.ParseSort(c.Query("sort"), nil)
	relation := dialect.ParseRelation(c.Query("expand"))

	rawPagination := dialect.ParsePagination(c.Query("pagination"))
	q := dialect.NewQuery(rawPagination, filter, sort, relation)

	sqlBoilerMods := q.Mods(dm.TableNames.Event)
	tychoSQL, tychoArgs := q.SQL(dm.TableNames.Event)

	records, err := dm.Events(sqlBoilerMods...).All(c, r.db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginatedRecords, pagination := query.Paginate(q, records)
	c.JSON(http.StatusOK, gin.H{
		"tychoSQL":         tychoSQL,
		"tychoArgs":        tychoArgs,
		"sqlBoilerMods":    sqlBoilerMods,
		"records":          records,
		"paginatedRecords": paginatedRecords,
		"pagination":       pagination,
	})
}

// GET /events/:id
func (r *Router) get(c *gin.Context) {
	relation := dialect.ParseRelation(c.Query("expand"))
	params := dialect.ParseParams(query.NewParam(dm.EventColumns.ID, c.Param("id")))
	q := dialect.NewQuery(relation, params)

  sqlBoilerMods := q.Mods(dm.TableNames.Event)
	tychoSQL, tychoArgs := q.SQL(dm.TableNames.Event)

	record, err := dm.Events(sqlBoilerMods...).One(c, r.db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
    "tychoSQL":         tychoSQL,
		"tychoArgs":        tychoArgs,
		"sqlBoilerMods":    sqlBoilerMods,
    "record": record,
  })
}

// GET /events/:id/comments
func (r *Router) listComments(c *gin.Context) {
	filter := dialect.ParseFilter(c.Query("filter"), nil)
	sort := dialect.ParseSort(c.Query("sort"), nil)
	relation := dialect.ParseRelation(c.Query("expand"))

	params := dialect.ParseParams(query.NewParam(dm.CommentColumns.EventID, c.Param("id")))

	rawPagination := dialect.ParsePagination(c.Query("pagination"))
	q := dialect.NewQuery(rawPagination, filter, sort, relation, params)

	sqlBoilerMods := q.Mods(dm.TableNames.Event)
	tychoSQL, tychoArgs := q.SQL(dm.TableNames.Event)

	records, err := dm.Comments(sqlBoilerMods...).All(c, r.db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginatedRecords, pagination := query.Paginate(q, records)
	c.JSON(http.StatusOK, gin.H{
		"tychoSQL":         tychoSQL,
		"tychoArgs":        tychoArgs,
		"sqlBoilerMods":    sqlBoilerMods,
		"records":          records,
		"paginatedRecords": paginatedRecords,
		"pagination":       pagination,
	})
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
in (in): []any
nin (not in): []any
c (contains): string
nc (not contains): string
sw (starts with): string
ew (ends with): string
null (is null): boolean
```

## Sorting
You can add multiple sorting columns. When the first one has duplicate values, it will sort by the next column etc.

```
https://domain.com/endpoint?sort=[{"colunn":"name", "order":"ASC"}]
```

## Relation
You can add the relations you want to include in the response.

```
https://domain.com/endpoint?relation=["table", "table2"]
```

## Param

``` 
https://domain.com/endpoint/:param
```


## Offset Pagination

```
https://domain.com/endpoint?pagination={"page": 1, "limit": 10}
```

## Cursor Pagination (backward does not work!)

```
https://domain.com/endpoint?pagination={"cursor": "optional cursor", "limit": 10}
```



## Typescript

```typescript
export interface Query {
  filter?: Filter;
  sort?: Sort[];
  pagination?: CursorPagination;
  expand?: string[];
  onBehalfOfAccountId?: string;
}

export type FilterType = 'eq' | 'neq' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'nin' | 'c' | 'nc' | 'sw' | 'ew' | 'null';
export type FilterTypeValue = string | number | boolean | string[];
export type FilterColumn = Record<string, Partial<Record<FilterType | 'or', FilterTypeValue>>>;
export type Filter = FilterColumn | Record<'or', FilterColumn>;

export type Sort = {
  column: string;
  order: 'ASC' | 'DESC';
};

export type CursorPagination = {
  limit: number;
  cursor?: string;
  page?: number;
};

export function createQuery(q: Query | undefined) {
  if (!q) {
    return '';
  }
  return Object.entries(q)
    .map(([key, value]) => {
      if (typeof value === 'object') {
        const newValue = removeEmptyTreeValues(value);
        if (newValue) {
          if (typeof newValue === 'object') {
            return `&${key}=${encodeURIComponent(JSON.stringify(newValue))}`;
          }

          return `&${key}=${encodeURIComponent(newValue)}`;
        }

        return false;
      }

      if (value !== undefined && value !== null) {
        return `&${key}=${encodeURIComponent(value)}`;
      }
    })
    .filter(Boolean)
    .join('')
    .replace('&', '?');
}

function removeEmptyTreeValues(obj: Record<string, any> | undefined): Record<string, any> | undefined {
  if (!obj) {
    return undefined;
  }

  if (Array.isArray(obj)) {
    if (obj.length === 0) {
      return undefined;
    }

    return obj.filter(Boolean);
  }

  const object: Record<string, any> = {};

  for (const [key, value] of Object.entries(obj)) {
    if (typeof value === 'object') {
      const treeValues = removeEmptyTreeValues(value);
      if (!treeValues) continue;
      if (Object.keys(treeValues).length > 0) {
        object[key] = treeValues;
      }
      continue;
    }

    if (value) {
      object[key] = value;
    }
  }

  return object;
}

```

## License

MIT © [Pointers](./LICENSE) 2024
