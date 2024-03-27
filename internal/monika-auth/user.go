package monikaauth

import (
	"fmt"

	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (auth *MonikaAuth) CreateUser(newUser *sdk.User) error {
	newUser.Username.Sanitize()

	if !newUser.Username.Valid() {
		return fmt.Errorf("username invalid: {length: 4-30, set: [a-z, A-Z, 0-9, `@`, `.`, `_`, `-`]}")
	}

	if !newUser.Role.Valid() {
		return fmt.Errorf("role invalid: unknown")
	}

	if !newUser.Password.Valid() {
		return fmt.Errorf("password invalid: {length 10-32, set: [all except whitespace]}")
	}

	if err := auth.hashPassword(newUser); err != nil {
		return err
	}

	return auth.userDB.CreateUser(newUser)
}

func (auth *MonikaAuth) GetAllUsers() ([]sdk.User, error) {
	return auth.userDB.GetAllUsers()
}

func (auth *MonikaAuth) GetUserById(userId int64) (*sdk.User, error) {
	user, err := auth.userDB.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	user.Password = "XXXXXXXX"
	return user, nil
}

func (auth *MonikaAuth) GetUserByUsername(userName sdk.Username) (*sdk.User, error) {
	user, err := auth.userDB.GetUserByUsername(userName)
	if err != nil {
		return nil, err
	}

	user.Password = "XXXXXXXX"
	return user, nil
}

func (auth *MonikaAuth) UpdateUserPassword(user *sdk.User) error {
	if !user.Password.Valid() {
		return fmt.Errorf("password invalid: {length 10-32, set: [all except whitespace]}")
	}

	// hash new password
	if err := auth.hashPassword(user); err != nil {
		return err
	}

	// update password in database
	return auth.userDB.UpdateUserPassword(user.UserId, user.Password)
}

func (auth *MonikaAuth) UpdateUserRole(user sdk.User) error {
	if !user.Role.Valid() {
		return fmt.Errorf("role invalid: unknown")
	}
	// update role in database
	return auth.userDB.UpdateUserRole(user.UserId, user.Role)
}

func (auth *MonikaAuth) DeleteUser(user sdk.User) error {
	return auth.userDB.DeleteUser(user.UserId)
}
