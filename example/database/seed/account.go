package seed

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/wearepointers/tycho/example/models/dm"
)

func Account(ctx context.Context, db *sql.DB) (*dm.Account, error) {
	d := &dm.Account{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@test.com",
		Password:  "$2a$10$aPuVYVEX6BOcW51e6sxWHOWCGtddQkAe6Q1AdlyUJ20f8d70Mi/aa", // Test123
	}

	if err := d.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	return d, nil
}
