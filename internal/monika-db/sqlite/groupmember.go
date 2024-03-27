package sqlite

import (
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (s *SQLiteRepository) CreateMember(module *sdk.GroupMember) error {
	if res, err := s.db.Exec(CREATE_MEMBER, module.ModuleName, module.ModuleType, module.DeviceName, module.DeviceType, module.Group); err != nil {
		return err
	} else {
		if id, err := res.LastInsertId(); err != nil {
			return err
		} else {
			module.Id = id
		}
	}

	return nil
}

func (s *SQLiteRepository) GetMember(memberId int64) (*sdk.GroupMember, error) {
	row := s.db.QueryRow(QUERY_MEMBER, memberId)

	member := sdk.GroupMember{}
	if err := row.Scan(&member.Id, &member.ModuleName, &member.ModuleType, &member.DeviceName, &member.DeviceType, &member.Group); err != nil {
		return nil, err
	}

	return &member, nil
}

func (s *SQLiteRepository) GetAllMembers() ([]sdk.GroupMember, error) {
	groupMembers := make([]sdk.GroupMember, 0)
	rows, err := s.db.Query(QUERY_ALL_MEMBERS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	member := sdk.GroupMember{}

	for rows.Next() {
		if err := rows.Scan(&member.Id, &member.ModuleName, &member.ModuleType, &member.DeviceName, &member.DeviceType, &member.Group); err != nil {
			return nil, err
		}

		groupMembers = append(groupMembers, member)
	}

	return groupMembers, nil
}

func (s *SQLiteRepository) GetMembersByGroup(group int64) ([]sdk.GroupMember, error) {
	groupMembers := make([]sdk.GroupMember, 0)
	rows, err := s.db.Query(QUERY_MEMBERS_BY_GROUP, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var member sdk.GroupMember

	for rows.Next() {
		if err := rows.Scan(&member.Id, &member.ModuleName, &member.ModuleType, &member.DeviceName, &member.DeviceType, &member.Group); err != nil {
			return nil, err
		}

		groupMembers = append(groupMembers, member)
	}

	return groupMembers, nil
}

func (s *SQLiteRepository) DeleteMember(moduleId int64) error {
	_, err := s.db.Exec(DELETE_MEMBER, moduleId)
	return err
}
