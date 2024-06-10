package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/wearepointers/tycho/example/models/dm"
)

func Events(ctx context.Context, db *sql.DB, organizationID, AccountID string) (dm.EventSlice, error) {
	var eventSlice dm.EventSlice
	for i := 0; i < 10000; i++ {
		f := &dm.Event{
			Name:           fmt.Sprintf("Event %d", i),
			Description:    fmt.Sprintf("Description %d", i),
			AccountID:      AccountID,
			OrganizationID: organizationID,
		}
		if err := f.Insert(ctx, db, boil.Infer()); err != nil {
			return nil, err
		}

		eventSlice = append(eventSlice, f)
		// Avoid same created_at
		time.Sleep(50 * time.Millisecond)
	}

	return eventSlice, nil
}
