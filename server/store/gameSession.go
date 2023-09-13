package store

import (
	"database/sql"
	"math/rand"
	"time"
)

type GameState int

const (
	Created GameState = iota
	Joining
	Played
)

const (
	CreatedTimeoutDuration time.Duration = time.Hour * 8
	JoiningTimeoutDuration time.Duration = time.Hour * 24
	PlayedTimeoutDuration  time.Duration = time.Hour * 72
)

type GameRuleSet int

const (
	Default GameRuleSet = iota
	WithBlacklist
)

type GameSession struct {
	Id      int64
	Created int64
	Key     string
	State   GameState
	RuleSet GameRuleSet

	Users *[]User
}

func FindSessionByID(id int64, db *sql.DB) (GameSession, error) {
	var session GameSession
	sql := `SELECT * FROM sessions WHERE id=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return session, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(
		&session.Id,
		&session.Created,
		&session.State,
		&session.Key,
		&session.RuleSet,
	)

	if err != nil {
		return session, err
	}

	users, err := FindUsersBySessionId(session.Id, db)
	if err == nil {
		session.Users = &users
	}

	return session, nil
}

func FindSessionByKey(key string, db *sql.DB) (GameSession, error) {
	var session GameSession
	sql := `SELECT * FROM sessions WHERE key=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return session, err
	}

	row := stmt.QueryRow(key)
	err = row.Scan(
		&session.Id,
		&session.Created,
		&session.State,
		&session.Key,
		&session.RuleSet,
	)

	if err != nil {
		return session, err
	}

	return session, nil
}

func CreateSession(db *sql.DB) (GameSession, error) {
	var session GameSession
	sql := `INSERT INTO sessions (created, state, key, rule_set) VALUES(?,?,?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return session, err
	}

	session.Created = time.Now().Unix()
	session.State = Created

	// fix potential collions at some point
	session.Key = create_random_unique_key()
	session.RuleSet = Default

	result, err := stmt.Exec(
		&session.Created,
		&session.State,
		&session.Key,
		&session.RuleSet,
	)
	if err != nil {
		return session, err
	}

	session.Id, err = result.LastInsertId()
	if err != nil {
		return session, err
	}

	return session, nil
}

func create_random_unique_key() string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomString := make([]byte, 16)
	for i := range randomString {
		randomString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randomString)
}

func (session *GameSession) Delete(db *sql.DB) error {
	sql := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(sql, session.Id)
	if err != nil {
		return err
	}
	return nil
}

func (session *GameSession) Update(db *sql.DB) error {
	sql := `
		UPDATE sessions SET state=? rule_set=? WHERE id=?
		WHERE users.id=?`
	stm, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stm.Exec(
		&session.State,
		&session.RuleSet,
		&session.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (session *GameSession) RollPartners(db *sql.DB) error {
	users, err := FindUsersBySessionId(session.Id, db)
	if err != nil {
		return err
	}

	for _, user := range users {
		user.PartnerId = 69
		user.Update(db)
	}

	return nil
}

func FindExpiredSessions(db *sql.DB) ([]GameSession, error) {
	var sessions []GameSession

	sql := `
		SELECT *
		FROM sessions
		WHERE (state = 0 AND created > ?)
		OR (state = 1 AND created > ?)
		OR (state = 2 AND created > ?)
	`

	now := time.Now()

	result, err := db.Query(
		sql,
		now.Add(-CreatedTimeoutDuration).Unix(),
		now.Add(-JoiningTimeoutDuration).Unix(),
		now.Add(-PlayedTimeoutDuration).Unix(),
	)
	if err != nil {
		return sessions, err
	}

	for result.Next() {
		var session GameSession
		err = result.Scan(
			&session.Id,
			&session.Created,
			&session.State,
			&session.Key,
			&session.RuleSet,
		)

		if err != nil {
			return sessions, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
