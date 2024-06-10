package database

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func migrater(reset bool) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Info().Str("dir", dir).Msg("Initializing database migration")
	m, err := migrate.New(
		fmt.Sprintf("file://%s/database/migrations", dir),
		getConnectionString(),
	)
	if err != nil {
		return err
	}

	if reset {
		log.Info().Msg("Resetting database")
		if err := m.Force(-1); err != nil {
			return err
		}
	}

	log.Info().Msg("Migrating database")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	err, err2 := m.Close()
	if err != nil {
		return err
	}

	if err2 != nil {
		return err
	}

	log.Info().Msg("Database migration complete")
	return nil
}
