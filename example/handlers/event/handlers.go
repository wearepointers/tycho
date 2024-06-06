package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wearepointers/tycho/example/models/dm"
	"github.com/wearepointers/tycho/example/models/erg"
	"github.com/wearepointers/tycho/query"
)

var dialect = query.Dialect{
	Driver:             query.Postgres,
	HasAutoIncrementID: false,
	APICasing:          query.CamelCase,
	DBCasing:           query.SnakeCase,
	PaginationType:     query.OffsetPagination,
	MaxLimit:           10,
}

func (r *Router) list(c *gin.Context) {
	filter := dialect.ParseFilter(c.Query("filter"), nil)
	sort := dialect.ParseSort(c.Query("sort"), nil)
	relation := dialect.ParseRelation(c.Query("expand"))

	rawPagination := dialect.ParsePagination(c.Query("pagination"))
	q := dialect.NewQuery(rawPagination, filter, sort, relation)

	records, err := dm.Events(q.Mods(table)...).All(c, r.db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginatedRecords, pagination := query.Paginate(q, records)
	c.JSON(http.StatusOK, gin.H{
		"records":    erg.ToEvents(paginatedRecords, nil),
		"pagination": pagination,
	})
}

func (r *Router) get(c *gin.Context) {
	relation := dialect.ParseRelation(c.Query("expand"))
	params := dialect.ParseParams(query.NewParam(dm.EventColumns.ID, c.Param("id")))
	q := dialect.NewQuery(relation, params)

	record, err := dm.Events(q.Mods(table)...).One(c, r.db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, erg.ToEvent(record, nil))
}

func (r *Router) listComments(c *gin.Context) {
	filter := dialect.ParseFilter(c.Query("filter"), nil)
	sort := dialect.ParseSort(c.Query("sort"), nil)
	relation := dialect.ParseRelation(c.Query("expand"))

	params := dialect.ParseParams(query.NewParam(dm.CommentColumns.EventID, c.Param("id")))

	rawPagination := dialect.ParsePagination(c.Query("pagination"))
	q := dialect.NewQuery(rawPagination, filter, sort, relation, params)

	// So now we need the param of events, but only based on event_id
	records, err := dm.Comments(q.Mods(dm.TableNames.Comment)...).All(c, r.db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginatedRecords, pagination := query.Paginate(q, records)
	c.JSON(http.StatusOK, gin.H{
		"records":    erg.ToComments(paginatedRecords, nil),
		"pagination": pagination,
	})
}
