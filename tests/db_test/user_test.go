package db_test

import (
	"testing"

	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func TestUsers(t *testing.T) {
	db, close := monikadb.NewSQLiteDatabase(":memory:", "admin")
	defer close()

	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}

	goodUser := sdk.User{
		Username: "Peter",
		Password: "geht dich nichts an",
		Role:     "admin",
	}

	anotherGoodUser := sdk.User{
		Username: "Charles",
		Password: "'1234?` SELECT",
		Role:     "admin",
	}

	badUser := sdk.User{
		Username: "Peter",
		Password: "nix",
		Role:     "admin",
	}

	if err := db.CreateUser(&goodUser); err != nil {
		t.Fatal(err)
	}

	if user, err := db.GetUserByUsername("Peter"); err != nil {
		t.Fatal(err)
	} else {
		if user.UserId != 2 {
			t.Logf("UserId wrong. Expected 2, got %d", user.UserId)
		}
	}

	if err := db.CreateUser(&anotherGoodUser); err != nil {
		t.Fatal(err)
	}

	if err := db.CreateUser(&badUser); err == nil {
		t.Fatal("should dbrt duplicate user")
	} else {
		t.Log(err)
	}

	if err := db.UpdateUserPassword(goodUser.UserId, "new password"); err != nil {
		t.Fatal(err)
	}

	if err := db.UpdateUserRole(goodUser.UserId, "normal"); err != nil {
		t.Fatal(err)
	}

	if returnedUsers, err := db.GetAllUsers(); err != nil {
		t.Fatal(err)
	} else {
		if len(returnedUsers) != 3 {
			t.Error("length of users does not match (should have 2)")
		}
		for _, users := range returnedUsers {
			t.Log(users)
		}
	}

	if err := db.DeleteUser(goodUser.UserId); err != nil {
		t.Fatal(err)
	}

	if returnedUsersAfterDelete, err := db.GetAllUsers(); err != nil {
		t.Log(err)
	} else {
		if len(returnedUsersAfterDelete) != 2 {
			t.Error("length of users does not match (should have 2)")
		}
	}
}
