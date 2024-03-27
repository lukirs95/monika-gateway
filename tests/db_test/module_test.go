package db_test

import (
	"fmt"
	"testing"

	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func TestMembers(t *testing.T) {
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
		t.Errorf("Wrong Group Id. Wanted: 1, got %d", goodGroup.Id)
	}

	Members := make([]sdk.GroupMember, 5000)
	for index := range Members {
		Members[index].ModuleName = fmt.Sprintf("Member %d", index)
		Members[index].ModuleType = sdk.ModuleType_AV
		Members[index].DeviceName = fmt.Sprintf("Device %d", index)
		Members[index].DeviceType = sdk.DeviceType_DIRECTOUT_RAVIO
		Members[index].Group = goodGroup.Id
	}

	for _, Member := range Members {
		if err := db.CreateMember(&Member); err != nil {
			t.Fatal(err)
		}
	}

	if allMembers, err := db.GetMembersByGroup(goodGroup.Id); err != nil {
		t.Fatal(err)
	} else if len(allMembers) != 5000 {
		t.Errorf("Wrong Members size. Wanted: 5000, got %d", len(allMembers))
	}

	if err := db.DeleteGroup(goodGroup.Id); err != nil {
		t.Error(err)
	}

	if allMembers, err := db.GetAllMembers(); err != nil {
		t.Fatal(err)
	} else if len(allMembers) != 0 {
		t.Errorf("Wrong Members size. Wanted: 0, got %d", len(allMembers))
	}

	if allGroups, err := db.GetAllGroups(); err != nil {
		t.Fatal(err)
	} else if len(allGroups) != 0 {
		t.Errorf("After deletion this should be 0, got: %d", len(allGroups))
	}

	if err := db.CreateGroup(&goodGroup); err != nil {
		t.Fatal(err)
	}

	if err := db.CreateMember(&Members[0]); err != nil {
		t.Fatal(err)
	}

	if allMembers, err := db.GetAllMembers(); err != nil {
		t.Fatal(err)
	} else if len(allMembers) != 1 {
		t.Errorf("Wrong Members size. Wanted: 1, got %d", len(allMembers))
	}

	if allGroups, err := db.GetAllGroups(); err != nil {
		t.Fatal(err)
	} else if len(allGroups) != 1 {
		t.Errorf("Wrong Members size. Wanted: 1, got %d", len(allGroups))
	}

	if err := db.DeleteMember(Members[0].Id); err != nil {
		t.Error(err)
	}

	if err := db.DeleteGroup(Members[0].Group); err != nil {
		t.Error(err)
	}

	if allMembers, err := db.GetAllMembers(); err != nil {
		t.Fatal(err)
	} else if len(allMembers) != 0 {
		t.Errorf("Wrong Members size. Wanted: 0, got %d", len(allMembers))
	}

	if allGroups, err := db.GetAllGroups(); err != nil {
		t.Fatal(err)
	} else if len(allGroups) != 0 {
		t.Errorf("After deletion this should be 0, got: %d", len(allGroups))
	}
}
