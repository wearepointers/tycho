package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/wearepointers/tycho/example/models/dm"
)

func Comments(ctx context.Context, db *sql.DB, eventID, organizationID, AccountID string) error {
	for i := 0; i < 100; i++ {
		f := &dm.Comment{
			Comment:        fmt.Sprintf("Comment %d", i),
			EventID:        eventID,
			AccountID:      AccountID,
			OrganizationID: organizationID,
		}
		if err := f.Insert(ctx, db, boil.Infer()); err != nil {
			return err
		}
		// Avoid same created_at
		time.Sleep(100 * time.Millisecond)

	}

	return nil
}
