package sqlite

import (
	"fmt"

	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (s *SQLiteRepository) CreateUser(user *sdk.User) error {
	if res, err := s.db.Exec(CREATE_USER, user.Username, user.Password, user.Role); err != nil {
		return err
	} else {
		if id, err := res.LastInsertId(); err != nil {
			return err
		} else {
			user.UserId = id
		}
	}

	return nil
}

func (s *SQLiteRepository) GetAllUsers() ([]sdk.User, error) {
	users := make([]sdk.User, 0)
	rows, err := s.db.Query(QUERY_ALL_USERS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user sdk.User
		if err := rows.Scan(&user.UserId, &user.Username, &user.Role); err != nil {
			return nil, err
		}

		user.Password = "XXXXXXXX"
		users = append(users, user)
	}
	return users, nil
}

func (s *SQLiteRepository) GetUserById(userId int64) (*sdk.User, error) {
	row := s.db.QueryRow(QUERY_USER_BY_ID, userId)

	user := sdk.User{}
	if err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SQLiteRepository) GetUserByUsername(username sdk.Username) (*sdk.User, error) {
	row := s.db.QueryRow(QUERY_USER_BY_USERNAME, username)

	user := sdk.User{}
	if err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SQLiteRepository) UpdateUserPassword(userId int64, password sdk.Password) error {
	_, err := s.db.Exec(UPDATE_USER_PASSWORD, password, userId)
	return err
}

func (s *SQLiteRepository) UpdateUserRole(userId int64, role sdk.UserRole) error {
	if userId == 1 {
		return fmt.Errorf("the admin role of superuser can't be changed")
	}
	_, err := s.db.Exec(UPDATE_USER_ROLE, role, userId)
	return err
}

func (s *SQLiteRepository) DeleteUser(userId int64) error {
	if userId == 1 {
		return fmt.Errorf("the superuser can't be deleted")
	}
	_, err := s.db.Exec(DELETE_USER, userId)
	return err
}
