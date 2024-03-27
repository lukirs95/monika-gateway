package db_test

import (
	"testing"

	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func TestGroups(t *testing.T) {
	db, close := monikadb.NewSQLiteDatabase(":memory:", "admin")
	defer close()

	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}

	goodGroup := sdk.Group{
		Name: "VenueKit 1",
	}

	if err := db.CreateGroup(&goodGroup); err != nil {
		t.Fatal(err)
	}

	if goodGroup.Id != 1 {
		t.Errorf("Wrong Group ID. Wanted: 1, got %d", goodGroup.Id)
	}

	if allGroups, err := db.GetAllGroups(); err != nil {
		t.Fatal(err)
	} else if len(allGroups) != 1 {
		t.Errorf("Wrong Groups size. Wanted: 1, got %d", len(allGroups))
	}

	if err := db.UpdateGroupname(goodGroup.Id, "Venue Kit 1"); err != nil {
		t.Error(err)
	}

	badGroup := sdk.Group{
		Name: "Venue Kit 1",
	}

	if err := db.CreateGroup(&badGroup); err == nil {
		t.Error("Should error because Groupname is already in table groups")
	}

	if err := db.DeleteGroup(goodGroup.Id); err != nil {
		t.Error(err)
	}

	if allGroups, err := db.GetAllGroups(); err != nil {
		t.Fatal(err)
	} else if len(allGroups) != 0 {
		t.Errorf("After deletion this should be 0, got: %d", len(allGroups))
	}
}
