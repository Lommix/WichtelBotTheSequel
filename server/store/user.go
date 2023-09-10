package store

import (
	"database/sql"
	"time"
)

type User struct {
	Id         int64
	Session_id int64
	Created    int64
	Name       string
	Password   string
	PartnerId  int64
	GroupId    string
	Notice     string
	Allergies  string
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
		&user.GroupId,
		&user.Notice,
		&user.Allergies,
	)
	if err != nil {
		return user, err
	}
	return user, nil
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
		&user.GroupId,
		&user.Notice,
		&user.Allergies,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateUser(user *User, db *sql.DB) error {
	sql := `INSERT INTO users (Session_id, Created, Name, Password, Notice, Allergies)
			VALUES (?,?,?,?,?,?);`
	stm, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	result, err := stm.Exec(user.Session_id, time.Now().Unix(), user.Name, user.Password, user.Notice, user.Allergies)
	if err != nil {
		return err
	}

	user.Id, _ = result.LastInsertId()

	return nil
}

func (user *User) Update(db *sql.DB) error {
	sql := `
		UPDATE users SET name=? password=? partner_id=? group_id=? notice=? allergies=?
		WHERE users.id=?`
	stm, err := db.Prepare(sql)

	if err != nil {
		return err
	}
	_, err = stm.Exec(
		&user.Name,
		&user.Password,
		&user.PartnerId,
		&user.GroupId,
		&user.Notice,
		&user.Allergies,
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
