package database

import (
	"context"
	"database/sql"

	"github.com/wearepointers/tycho/example/database/seed"
)

func seeder(db *sql.DB) error {
	ctx := context.Background()

	a, err := seed.Account(ctx, db)
	if err != nil {
		return err
	}

	o, err := seed.Organization(ctx, db)
	if err != nil {
		return err
	}

	if _, err := seed.OrganizationAccount(ctx, db, o.ID, a.ID); err != nil {
		return err
	}

	if err := seed.Events(ctx, db, o.ID, a.ID); err != nil {
		return err
	}

	return nil
}
