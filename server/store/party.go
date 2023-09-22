package store

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

const (
	CreatedTimeoutDuration time.Duration = time.Hour * 24
	JoiningTimeoutDuration time.Duration = time.Hour * 72
	PlayedTimeoutDuration  time.Duration = time.Hour * 72
)

type GameState int

const (
	Created GameState = iota
	Joining
	Played
)

type Party struct {
	Id        int64
	Created   int64
	Key       string
	State     GameState
	Blacklist bool
	Users     *[]User
}

func FindPartyByID(id int64, db *sql.DB) (Party, error) {
	var party Party
	sql := `SELECT * FROM parties WHERE id=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return party, err
	}

	row := stmt.QueryRow(id)
	err = row.Scan(
		&party.Id,
		&party.Created,
		&party.State,
		&party.Key,
		&party.Blacklist,
	)

	if err != nil {
		return party, err
	}

	users, err := FindUsersByPartyId(db, party.Id)
	if err == nil {
		party.Users = &users
	}

	return party, nil
}

func FindPartyByKey(key string, db *sql.DB) (Party, error) {
	var party Party
	sql := `SELECT * FROM parties WHERE key=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return party, err
	}

	row := stmt.QueryRow(key)
	err = row.Scan(
		&party.Id,
		&party.Created,
		&party.State,
		&party.Key,
		&party.Blacklist,
	)

	if err != nil {
		return party, err
	}

	user, err := FindUsersByPartyId(db, party.Id)
	if err == nil {
		party.Users = &user
	}

	return party, nil
}

func CreateParty(db *sql.DB) (Party, error) {
	var party Party
	sql := `INSERT INTO parties (created, state, key) VALUES(?,?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return party, err
	}

	party.Created = time.Now().Unix()
	party.State = Created

	party.Key, err = createRandomUniqueKey(db)
	if err != nil {
		return party, err
	}

	result, err := stmt.Exec(
		&party.Created,
		&party.State,
		&party.Key,
	)
	if err != nil {
		return party, err
	}

	party.Id, err = result.LastInsertId()
	if err != nil {
		return party, err
	}

	return party, nil
}

func (party *Party) Delete(db *sql.DB) error {
	sql := `DELETE FROM parties WHERE id = ?`
	_, err := db.Exec(sql, party.Id)
	if err != nil {
		return err
	}
	return nil
}

func (party *Party) Update(db *sql.DB) error {
	sql := `UPDATE parties SET state=?, blacklist=? WHERE id=?`
	stm, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stm.Exec(
		&party.State,
		&party.Blacklist,
		&party.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (party *Party) RollPartners(db *sql.DB) error {
	users, err := FindUsersByPartyId(db, party.Id)
	if err != nil {
		return err
	}

	if len(users) < 2 {
		return errors.New("Not enough players in party")
	}

	availablePartners := make([]User, len(users))
	copy(availablePartners, users)

	for _, user := range users {
		var indexList []int
		for i, u := range availablePartners {
			if party.Blacklist && user.ExcludeId == u.Id {
				continue
			}
			if user.Id != u.Id {
				indexList = append(indexList, i)
			}
		}

		if len(indexList) == 0 {
			return errors.New("No potential partners found, did you all blacklist the same person?")
		}

		partnerIndex := indexList[rand.Intn(len(indexList))]
		user.PartnerId = availablePartners[partnerIndex].Id
		user.Update(db)

		availablePartners = append(availablePartners[:partnerIndex], availablePartners[partnerIndex+1:]...)
	}

	party.State = Played
	return party.Update(db)
}

func filter(slice []interface{}, fn func(interface{}) bool) []interface{} {
	var out []any
	for _, item := range slice {
		if fn(item) {
			out = append(out, item)
		}
	}
	return out
}

func FindExpiredParties(db *sql.DB) ([]Party, error) {
	var parties []Party

	sql := `
		SELECT *
		FROM parties
		WHERE (state = 0 AND created < ?)
		OR (state = 1 AND created < ?)
		OR (state = 2 AND created < ?)`

	now := time.Now()
	result, err := db.Query(
		sql,
		now.Add(-CreatedTimeoutDuration).Unix(),
		now.Add(-JoiningTimeoutDuration).Unix(),
		now.Add(-PlayedTimeoutDuration).Unix(),
	)

	if err != nil {
		return parties, err
	}

	for result.Next() {
		var party Party
		err = result.Scan(
			&party.Id,
			&party.Created,
			&party.State,
			&party.Key,
			&party.Blacklist,
		)

		if err != nil {
			return parties, err
		}

		parties = append(parties, party)
	}

	return parties, nil
}

func partyKeyExists(db *sql.DB, key string) bool {
	sql := `SELECT * FROM parties WHERE key=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return true
	}
	_, err = stmt.Query(key)
	if err != nil {
		return true
	}
	return false
}

func createRandomUniqueKey(db *sql.DB) (string, error) {
	key := createRandomKey()
	timeout := 0
	for partyKeyExists(db, key) {
		key = createRandomKey()
		timeout++
		// practical engineering, math is hard
		if timeout > 50 {
			return "", errors.New("Server are busy")
		}
	}

	return key, nil
}

func createRandomKey() string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomString := make([]byte, 16)
	for i := range randomString {
		randomString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randomString)
}
