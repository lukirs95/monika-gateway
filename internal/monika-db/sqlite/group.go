package sqlite

import (
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (s *SQLiteRepository) CreateGroup(group *sdk.Group) error {
	if res, err := s.db.Exec(CREATE_GROUP, group.Name); err != nil {
		return err
	} else {
		if id, err := res.LastInsertId(); err != nil {
			return err
		} else {
			group.Id = id
		}
	}

	return nil
}

func (s *SQLiteRepository) GetGroup(groupId int64) (*sdk.Group, error) {
	group := sdk.Group{}
	row := s.db.QueryRow(QUERY_GROUP, groupId)

	if err := row.Scan(&group.Id, &group.Name); err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *SQLiteRepository) GetAllGroups() ([]sdk.Group, error) {
	groups := make([]sdk.Group, 0)
	rows, err := s.db.Query(QUERY_ALL_GROUPS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	group := sdk.Group{}

	for rows.Next() {
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (s *SQLiteRepository) UpdateGroupname(id int64, newGroupname string) error {
	_, err := s.db.Exec(UPDATE_GROUPNAME, newGroupname, id)
	return err
}

func (s *SQLiteRepository) DeleteGroup(id int64) error {
	_, err := s.db.Exec(DELETE_GROUP, id, id)
	return err
}
