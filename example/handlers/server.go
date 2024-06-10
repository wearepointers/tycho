package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wearepointers/tycho/example/handlers/event"
)

func Run(db *sql.DB, addr string) error {
	log.Info().Msg("Initializing server")

	router := gin.Default()

	v1Group := router.Group("/v1")
	event.Routes(v1Group, db)

	log.Info().Msgf("Server listening on %v", addr)
	return router.Run(addr)
}
