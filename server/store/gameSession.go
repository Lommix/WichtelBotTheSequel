package store

import (
	"database/sql"
	"math/rand"
	"time"
)

type GameState int
const(
	Created GameState = iota
	Joining
	Played
)

type GameRuleSet int
const(
	// the basic game
	Default GameRuleSet = iota
	// couples can exclude each other
	CoupleBlacklist
)

type GameSession struct{
	Id int64
	Created int64
	Key string
	State GameState
	Rules string
}


func FindSessionByID(id int64, db *sql.DB) (GameSession, error){
	var session GameSession
	sql := `SELECT * FROM sessions WHERE id=?`
	stmt, err := db.Prepare(sql)
	if err!= nil{
		return session, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(
		&session.Id,
		&session.Created,
		&session.State,
		&session.Key,
		&session.Rules,
	)

	if err!= nil{
		return session, err
	}

	return session, nil
}


func FindSessionByKey(key string, db *sql.DB) (GameSession, error){
	var session GameSession
	sql := `SELECT * FROM sessions WHERE key=?`
	stmt, err := db.Prepare(sql)
	if err!= nil{
		return session, err
	}

	row := stmt.QueryRow(key)
	err = row.Scan(
		&session.Id,
		&session.Created,
		&session.State,
		&session.Key,
		&session.Rules,
	)

	if err!= nil{
		return session, err
	}

	return session, nil
}

func CreateSession(db *sql.DB,) (GameSession, error){
	var session GameSession
	sql := `INSERT INTO session (created, state, key, rules) VALUES(?,?,?,?)`
	stmt, err := db.Prepare(sql)
	if err!= nil{
		return session, err
	}

	session.Created = time.Now().Unix()
	session.State = Created

	//fix potential collions at some point
	session.Key = create_random_unique_key()
	session.Rules = "none"

	result, err := stmt.Exec(
		&session.Created,
		&session.State,
		&session.Key,
		&session.Rules,
	)

	if err!= nil{
		return session, err
	}

	session.Id, err = result.LastInsertId()
	if err!= nil{
		return session, err
	}

	return session, nil
}


func create_random_unique_key() string{
	rand.Seed(time.Now().UnixNano())
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
