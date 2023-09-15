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

	Party *Party
	Partner     *User
}

func FindUserById(db *sql.DB, id int64) (User, error) {
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
		user.Party = &session
	}


	return user, nil
}

func FindUsersByPartyId(db *sql.DB, id int64) ([]User, error) {
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

func FindUserByNameAndRoomKey(db *sql.DB, name string, roomKey string) (User, error) {
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

// TODO: full switch
// optimizing db calls for one query per request
func FindUserWithPartyFast(db *sql.DB, userId int64) (User, error) {
	var requestedUser User
	var members []User
	var party Party
	sql := `SELECT * FROM users
			JOIN parties ON users.party_id = parties.id
			WHERE party_id = (SELECT party_id FROM users WHERE id = ?)`

	result, err := db.Query(sql, userId)
	if err != nil {
		return requestedUser, err
	}

	for result.Next(){
		var u User
		err := result.Scan(
			&u.Id,
			&u.PartyId,
			&u.Created,
			&u.Name,
			&u.Password,
			&u.PartnerId,
			&u.ExcludeId,
			&u.Notice,
			&u.Role,
			&party.Id,
			&party.Created,
			&party.State,
			&party.Key,
			&party.RuleSet,
		)

		if err != nil {
			return u, err
		}

		u.Party = &party
		members = append(members, u)
	}

	party.Users = &members

	// find requested user
	for _, user := range members {
		if user.Id == userId {
			requestedUser = user
		}
	}

	// find requested user partner
	if requestedUser.PartyId != 0 {
		for _, user := range members {
			if user.Id == requestedUser.PartnerId {
				requestedUser.Partner = &user
			}
		}
	}

	return requestedUser, nil
}
