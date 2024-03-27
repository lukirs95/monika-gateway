package sqlite

// ------ USERS ------ //

// Check if table users exists in database.
const EXISTS_TABLE_USERS = `SELECT name FROM sqlite_master WHERE type='table' AND name='users';`

// Create table users.
const CREATE_TABLE_USERS = `CREATE TABLE users (
	user_id INTEGER PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	role TEXT NOT NULL
);`

// Create new user in table users. {username: string, password: string, role: string}
const CREATE_USER = `INSERT INTO users(user_id, username, password, role) values(NULL, ?, ?, ?);`

// Query all users in table users.
const QUERY_ALL_USERS = `SELECT user_id, username, role FROM users;`

// Query one user by user_id. {user_id: int64}
const QUERY_USER_BY_ID = `SELECT * FROM users WHERE user_id = ?;`

// Query one user by username from table users. {username: string}
const QUERY_USER_BY_USERNAME = `SELECT * FROM users WHERE username = ?;`

// Update password at given user in table users. {password: string, user_id: int64}
const UPDATE_USER_PASSWORD = `UPDATE users SET password = ? WHERE user_id = ?;`

// Update role at given user in table users. {role: string, user_id: int64}
const UPDATE_USER_ROLE = `UPDATE users SET role = ? WHERE user_id = ?;`

// Delete user from table users. {user_id: int64}
const DELETE_USER = `DELETE FROM users WHERE user_id = ?;`

// ------ END USERS ----- //

// ------ GROUPS ------ //

// Check if table groups exists in database.
const EXISTS_TABLE_GROUPS = `SELECT name FROM sqlite_master WHERE type='table' AND name='groups';`

// Create table groups.
const CREATE_TABLE_GROUPS = `CREATE TABLE groups (
  group_id INTEGER PRIMARY KEY,
  groupname TEXT NOT NULL UNIQUE
);`

// Create new group in table groups. {groupname: string}
const CREATE_GROUP = `INSERT INTO groups(group_id, groupname) values(NULL, ?);`

// Query group in table groups with groupid. {group_id: int64}
const QUERY_GROUP = `SELECT * FROM groups WHERE group_id = ?;`

// Query all groups in table groups.
const QUERY_ALL_GROUPS = `SELECT * FROM groups;`

// Update groupname in table groups. {groupname: string, group_id: int64}
const UPDATE_GROUPNAME = `UPDATE groups SET groupname = ? WHERE group_id = ?;`

// Delete group from table groups. Delete modules where group_id is group {group_id: int64, group_id: int64}
const DELETE_GROUP = `DELETE FROM members WHERE group_id = ?;
	DELETE FROM groups WHERE group_id = ?;`

// ------ END GROUPS ------ //

// ------ MEMBER ------ //

// Check if table groups exists in database.
const EXISTS_TABLE_MEMBERS = `SELECT name FROM sqlite_master WHERE type='table' AND name='members';`

// Create table members.
const CREATE_TABLE_MEMBERS = `CREATE TABLE members (
	member_id INTEGER PRIMARY KEY,
	modulename TEXT NOT NULL,
  moduletype TEXT NOT NULL,
	devicename TEXT NOT NULL,
  devicetype TEXT NOT NULL,
	group_id INTEGER NOT NULL,
	FOREIGN KEY(group_id) REFERENCES groups(group_id)
);`

// Create new member in table members. {modulename: string, moduletype: string, devicename: string, devicetype: string, group_id: int64}
const CREATE_MEMBER = `INSERT INTO members(
	member_id,
	modulename,
  moduletype,
	devicename,
  devicetype,
	group_id
) values(NULL, ?, ?, ?, ?, ?);`

// Query one member from table members by id. {member_id: int64}
const QUERY_MEMBER = `SELECT * FROM members WHERE member_id = ?;`

// Query all members from table members.
const QUERY_ALL_MEMBERS = `SELECT * FROM members ORDER BY group_id, devicename, modulename;`

// Query members by group from table members. {group_id: int64}
const QUERY_MEMBERS_BY_GROUP = `SELECT * FROM members WHERE group_id = ? ORDER BY devicetype, moduletype, devicename, modulename;`

// Delete module by module_id from table members. {member_id: int64}
const DELETE_MEMBER = `DELETE FROM members WHERE member_id = ?`

// ------ END MODULES ------ //
