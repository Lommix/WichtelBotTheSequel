package store

import (
	"database/sql"
	"time"
)

type UserRole int

const (
	DefaultUser UserRole = iota
	Moderator
)

type User struct {
	Id         int64
	Session_id int64
	Created    int64
	Name       string
	Password   string
	PartnerId  int64
	ExcludeId  int64
	Notice     string
	Allergies  string
	Role       UserRole
}

func FindUserById(id int, db *sql.DB) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(
		&user.Id,
		&user.Session_id,
		&user.Created,
		&user.Name,
		&user.Password,
		&user.PartnerId,
		&user.ExcludeId,
		&user.Notice,
		&user.Allergies,
		&user.Role,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func FindUsersBySessionId(id int64, db *sql.DB) ([]User, error) {
	var users []User
	sql := `SELECT * FROM users WHERE session_id=?`
	result, err := db.Query(sql, id)
	if err != nil {
		return users, err
	}
	for result.Next() {
		var user User
		err = result.Scan(
			&user.Id,
			&user.Session_id,
			&user.Created,
			&user.Name,
			&user.Password,
			&user.PartnerId,
			&user.ExcludeId,
			&user.Notice,
			&user.Allergies,
			&user.Role,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func FindUserByNameAndRoomKey(name string, roomKey string, db *sql.DB) (User, error) {
	var user User
	query := `SELECT * FROM users
			  JOIN sessions ON users.session_id = sessions.id
			  WHERE users.name = ? AND sessions.key = ?;`
	row := db.QueryRow(query, name, roomKey)
	err := row.Scan(
		&user.Id,
		&user.Session_id,
		&user.Created,
		&user.Name,
		&user.Password,
		&user.PartnerId,
		&user.ExcludeId,
		&user.Notice,
		&user.Allergies,
		&user.Role,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateUser(
	db *sql.DB,
	session_id int64,
	name string,
	password string,
	notice string,
	allergies string,
	role UserRole,
) (User, error) {
	var user User
	sql := `INSERT INTO users (session_id, created, name, password, notice, allergies, role)
			VALUES (?,?,?,?,?,?,?);`
	stm, err := db.Prepare(sql)
	if err != nil {
		return user, err
	}

	result, err := stm.Exec(
		session_id,
		time.Now().Unix(),
		name,
		password,
		notice,
		allergies,
		role,
	)
	if err != nil {
		return user, err
	}

	user.Id, _ = result.LastInsertId()
	user.Session_id = session_id
	user.Name = name
	user.Password = password
	user.Notice = notice
	user.Allergies = allergies
	user.Role = role

	return user, nil
}

func (user *User) Update(db *sql.DB) error {
	sql := `
		UPDATE users SET name=? password=? partner_id=? exclude_id=? notice=? allergies=? role=?
		WHERE users.id=?`
	stm, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stm.Exec(
		&user.Name,
		&user.Password,
		&user.PartnerId,
		&user.ExcludeId,
		&user.Notice,
		&user.Allergies,
		&user.Role,
		&user.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (user *User) Delete(db *sql.DB) error {
	sql := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(sql, user.Id)
	if err != nil {
		return err
	}
	return nil
}


func DeleteAllUserInSession(db *sql.DB, sessionId int64) error {
	sql := `DELETE FROM users WHERE sesssion_id=?`
	_, err := db.Exec(sql, sessionId)
	if err != nil {
		return err
	}
	return nil
}
