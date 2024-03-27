package dbrepo

import sdk "github.com/lukirs95/monika-gosdk/pkg/types"

type DatabaseRepo interface {
	UserDatabaseRepo
	GroupDatabaseRepo
}

type UserDatabaseRepo interface {
	Migrate() error
	CreateUser(user *sdk.User) error
	GetAllUsers() (users []sdk.User, err error)
	GetUserById(userId int64) (*sdk.User, error)
	GetUserByUsername(username sdk.Username) (*sdk.User, error)
	UpdateUserPassword(userId int64, password sdk.Password) error
	UpdateUserRole(userId int64, role sdk.UserRole) error
	DeleteUser(userId int64) error
}

type GroupDatabaseRepo interface {
	Migrate() error
	CreateGroup(group *sdk.Group) error
	GetGroup(groupId int64) (*sdk.Group, error)
	GetAllGroups() ([]sdk.Group, error)
	UpdateGroupname(id int64, newGroupname string) error
	DeleteGroup(groupId int64) error
	CreateMember(groupMember *sdk.GroupMember) error
	GetMember(memberId int64) (*sdk.GroupMember, error)
	GetAllMembers() (groupMembers []sdk.GroupMember, err error)
	GetMembersByGroup(group int64) (groupMembers []sdk.GroupMember, err error)
	DeleteMember(memberId int64) error
}
