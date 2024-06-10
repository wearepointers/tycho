package seed

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/wearepointers/tycho/example/models/dm"
)

func Organization(ctx context.Context, db *sql.DB) (*dm.Organization, error) {
	d := &dm.Organization{
		Name: "Organization 1",
	}

	if err := d.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	return d, nil
}

func OrganizationAccount(ctx context.Context, db *sql.DB, organizationID, accountID string) (*dm.OrganizationAccount, error) {
	d := &dm.OrganizationAccount{
		OrganizationID: organizationID,
		AccountID:      accountID,
		Role:           dm.RoleADMIN,
	}

	if err := d.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	return d, nil
}
