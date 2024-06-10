package event

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/wearepointers/tycho/example/models/dm"
)

var (
	table = dm.TableNames.Event
)

type Router struct {
	db *sql.DB
}

func Routes(r *gin.RouterGroup, db *sql.DB) {
	s := Router{db}

	r.GET("/event", s.list)
	r.GET("/event/:id", s.get)
	r.GET("/event/:id/comments", s.listComments)
}
