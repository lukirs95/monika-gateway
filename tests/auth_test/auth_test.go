package authtest

import (
	"testing"
	"time"

	monikaauth "github.com/lukirs95/monika-gateway/internal/monika-auth"
	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

func TestAuth(t *testing.T) {
	db, close := monikadb.NewSQLiteDatabase(":memory:", "admin")
	defer close()

	db.Migrate()

	auth := monikaauth.NewMonikaAuth(db, monikaauth.NewMonikaAuthParams(time.Second*1, []byte("HASH")))

	newUser := &types.User{
		Username: "correctUsername",
		Password: "correctPassword",
		Role:     types.UserRole_ADMIN,
	}

	if err := auth.CreateUser(newUser); err != nil {
		t.Fatal(err)
	}

	if newUser.Password == "correctPassword" {
		t.Fatal("the provided password has to be hashed")
	}

	newAuth := types.Auth{
		Username: "correctUsername",
		Password: "correctPassword",
		Expires:  true,
	}

	token, err := auth.Login(newAuth)
	if err != nil {
		t.Fatal(err)
	}

	session, err := auth.ValidateJWT(token)
	if err != nil {
		t.Fatal(err)
	}

	if session.UserId != 2 {
		t.Errorf("invalid session id. Wanted 2, got %d", session.UserId)
	}

	if session.Role != types.UserRole_ADMIN {
		t.Errorf("invalid session role. Wanted `admin`, got %s", session.Role)
	}

	newAuth.Password = "wrongPassword"

	_, err = auth.Login(newAuth)
	if err == nil {
		t.Errorf("wrong password got entry")
	}

	newAuth.Password = "correctPassword"
	newAuth.Username = "wrongUsername"

	_, err = auth.Login(newAuth)
	if err == nil {
		t.Errorf("wrong username got access")
	}

	err = auth.DeleteUser(*newUser)
	if err != nil {
		t.Error(err)
	}

	newAuth.Username = "correctUsername"

	_, err = auth.Login(newAuth)
	if err == nil {
		t.Errorf("user has not been deleted")
	}

	time.Sleep(time.Second * 2)
	_, err = auth.ValidateJWT(token)
	if err == nil {
		t.Error("session has not expired")
	}

}
