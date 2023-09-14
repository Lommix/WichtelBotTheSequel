package store

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

const (
	CreatedTimeoutDuration time.Duration = time.Hour * 8
	JoiningTimeoutDuration time.Duration = time.Hour * 24
	PlayedTimeoutDuration  time.Duration = time.Hour * 72
)

type GameState int
const (
	Created GameState = iota
	Joining
	Played
)

type GameRuleSet int
const (
	Default GameRuleSet = iota
	WithBlacklist
)

type Party struct {
	Id      int64
	Created int64
	Key     string
	State   GameState
	RuleSet GameRuleSet

	Users *[]User
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
		&party.RuleSet,
	)

	if err != nil {
		return party, err
	}

	users, err := FindUsersByPartyId(party.Id, db)
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
		&party.RuleSet,
	)

	if err != nil {
		return party, err
	}

	return party, nil
}

func CreateParty(db *sql.DB) (Party, error) {
	var party Party
	sql := `INSERT INTO parties (created, state, key, rule_set) VALUES(?,?,?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return party, err
	}

	party.Created = time.Now().Unix()
	party.State = Created

	var key string
	// timeout := 0
	// for {
	key = createRandomKey()
	// 	_, err := FindPartyByKey(key, db)
	// 	if err != nil {
	// 		timeout ++
	// 		continue
	// 	}
	//
	// 	if timeout > 50 {
	// 		return party, errors.New("Server are busy")
	// 	}
	// 	break
	// }

	party.Key = key
	party.RuleSet = Default

	result, err := stmt.Exec(
		&party.Created,
		&party.State,
		&party.Key,
		&party.RuleSet,
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
	sql := `UPDATE parties SET state=?, rule_set=? WHERE id=?`
	stm, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stm.Exec(
		&party.State,
		&party.RuleSet,
		&party.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (party *Party) RollPartners(db *sql.DB) error {
	users, err := FindUsersByPartyId(party.Id, db)
	if err != nil {
		return err
	}

	if len(users) < 2 {
		return errors.New("Not enough players in party")
	}

	for _, user := range users {

		var potentialPartner []User

//&& (user.ExcludeId != u.Id && party.RuleSet == WithBlacklist)

		for _, u := range users {
			if user.Id != u.Id && u.PartnerId == 0 {
				potentialPartner = append(potentialPartner, u)
			}
		}

		if len(potentialPartner) == 0 {
			return errors.New("No potential partners found, did you all blacklist the same person?")
		}

		partner := potentialPartner[rand.Intn(len(potentialPartner))]

		user.PartnerId = partner.Id
		user.Update(db)
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
		return parties, err
	}

	for result.Next() {
		var party Party
		err = result.Scan(
			&party.Id,
			&party.Created,
			&party.State,
			&party.Key,
			&party.RuleSet,
		)

		if err != nil {
			return parties, err
		}

		parties = append(parties, party)
	}

	return parties, nil
}

func createRandomKey() string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomString := make([]byte, 16)
	for i := range randomString {
		randomString[i] = chars[rand.Intn(len(chars))]
	}
	return string(randomString)
}
