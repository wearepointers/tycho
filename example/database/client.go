package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	ErrMigrated = errors.New("migration: success")
)

func New(migrate, reset, seed bool) (*sql.DB, error) {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		return nil, err
	}

	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(25)
	// db.SetConnMaxLifetime(time.Hour)

	if migrate {
		if err := migrater(reset); err != nil {
			return nil, err
		}

		if seed {
			if err := seeder(db); err != nil {
				return nil, err
			}
		}

		return db, ErrMigrated
	}

	return db, nil
}

func getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%v", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"))
}
