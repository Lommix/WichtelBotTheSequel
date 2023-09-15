package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Stats struct {
	id              int64
	created         time.Time
	games_played    int64
	user_registered int64
}

func GetDailyStats(day time.Time, db *sql.DB) (Stats, error) {
	var stats Stats
	stats.created = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	row := db.QueryRow(`SELECT * FROM stats WHERE created=?`, stats.created.Unix())
	err := row.Scan(
		&stats.id,
		&stats.created,
		&stats.games_played,
		&stats.user_registered,
	)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func AddGamePlayed(db *sql.DB) error {
	stats, err := FindOrCreateStats(db)
	if err != nil {
		return err
	}

	stats.games_played = stats.games_played + 1
	err = stats.Update(db)
	if err != nil {
		return err
	}

	return nil
}

func AddUserRegistered(db *sql.DB) error {
	stats, err := FindOrCreateStats(db)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	stats.user_registered = stats.user_registered + 1
	err = stats.Update(db)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func FindOrCreateStats(db *sql.DB) (Stats, error) {
	var stats Stats
	now := time.Now()
	stats.created = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	row := db.QueryRow(`SELECT id, games_played, user_registered FROM stats WHERE created = ?`, stats.created.Unix())

	err := row.Scan(
		&stats.id,
		&stats.games_played,
		&stats.user_registered,
	)

	if err != nil {
		println(err.Error())
		result, err := db.Exec(`INSERT INTO stats (created, games_played, user_registered) VALUES(?,0,0)`, stats.created.Unix())
		if err != nil {
			return stats, err
		}

		stats.id, err = result.LastInsertId()
		if err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func (stats *Stats) Update(db *sql.DB) error {
	sql := `UPDATE stats SET games_played=?, user_registered=? WHERE id=?`
	_, err := db.Exec(sql, &stats.games_played, &stats.user_registered, &stats.id)
	if err != nil {
		return err
	}
	return nil
}
