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
	Id        int64
	PartyId   int64
	Created   int64
	Name      string
	Password  string
	PartnerId int64
	ExcludeId int64
	Notice    string
	Role      UserRole

	GameSession *Party
	Partner     *User
}

func FindUserById(id int64, db *sql.DB) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(
		&user.Id,
		&user.PartyId,
		&user.Created,
		&user.Name,
		&user.Password,
		&user.PartnerId,
		&user.ExcludeId,
		&user.Notice,
		&user.Role,
	)
	if err != nil {
		return user, err
	}

	session, err := FindPartyByID(user.PartyId, db)
	if err == nil {
		user.GameSession = &session
	}


	return user, nil
}

func FindUsersByPartyId(id int64, db *sql.DB) ([]User, error) {
	var users []User
	sql := `SELECT * FROM users WHERE party_id=?`
	result, err := db.Query(sql, id)
	if err != nil {
		return users, err
	}
	for result.Next() {
		var user User
		err = result.Scan(
			&user.Id,
			&user.PartyId,
			&user.Created,
			&user.Name,
			&user.Password,
			&user.PartnerId,
			&user.ExcludeId,
			&user.Notice,
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
	query := `SELECT users.* FROM users
			  INNER JOIN parties ON users.party_id = parties.id
			  WHERE users.name = ? AND parties.key = ?;`
	row := db.QueryRow(query, name, roomKey)
	err := row.Scan(
		&user.Id,
		&user.PartyId,
		&user.Created,
		&user.Name,
		&user.Password,
		&user.PartnerId,
		&user.ExcludeId,
		&user.Notice,
		&user.Role,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateUser(
	db *sql.DB,
	partyId int64,
	name string,
	password string,
	notice string,
	role UserRole,
) (User, error) {
	var user User
	sql := `INSERT INTO users (party_id, created, name, password, notice, role)
			VALUES (?,?,?,?,?,?);`
	stm, err := db.Prepare(sql)
	if err != nil {
		return user, err
	}

	result, err := stm.Exec(
		partyId,
		time.Now().Unix(),
		name,
		password,
		notice,
		role,
	)
	if err != nil {
		return user, err
	}

	user.Id, _ = result.LastInsertId()
	user.PartyId = partyId
	user.Name = name
	user.Password = password
	user.Notice = notice
	user.Role = role

	return user, nil
}

func (user *User) Update(db *sql.DB) error {
	sql := `
		UPDATE users
		SET password = ?, partner_id = ?, exclude_id = ?, notice = ?, role = ?
		WHERE id=?`
	stm, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stm.Exec(
		&user.Password,
		&user.PartnerId,
		&user.ExcludeId,
		&user.Notice,
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

func DeleteUsersInParty(db *sql.DB, sessionId int64) error {
	sql := `DELETE FROM users WHERE party_id=?`
	_, err := db.Exec(sql, sessionId)
	if err != nil {
		return err
	}
	return nil
}
