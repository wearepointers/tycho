package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/wearepointers/tycho/example/database"
	"github.com/wearepointers/tycho/example/handlers"
)

var (
	dbMigrateFlag = flag.Bool("db:migrate", false, "run database migrations")
	dbSeedFlag    = flag.Bool("db:seed", false, "run database seeders")
	dbResetFlag   = flag.Bool("db:reset", false, "reset database migrations")
)

func main() {
	flag.Parse()
	godotenv.Load()

	var addr = fmt.Sprintf("%v:%v", os.Getenv("HOST"), os.Getenv("PORT"))

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	gin.SetMode(gin.DebugMode)
	boil.DebugMode = true

	db, err := database.New(*dbMigrateFlag, *dbResetFlag, *dbSeedFlag)
	if err != nil {
		if err != database.ErrMigrated {
			log.Fatal().Err(err).Msg("Failed to initialize database")
			return
		}

		return
	}
	defer db.Close()

	log.Fatal().Err(handlers.Run(db, addr)).Msgf("Failed to start API service")
}
