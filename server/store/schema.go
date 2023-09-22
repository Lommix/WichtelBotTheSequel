package store

import "database/sql"

func SchemaUp(db *sql.DB) error {
	sql := `
		CREATE TABLE IF NOT EXISTS parties (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			created INTEGER NOT NULL,
			state INTEGER NOT NULL,
			key STRING NOT NULL,
			blacklist BOOLEAN DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			party_id INTEGER NOT NULL,
			created INTEGER NOT NULL,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			partner_id INTEGER DEFAULT 0,
			exclude_id TEXT DEFAULT 0,
			notice TEXT DEFAULT NULL,
			role INTEGER DEFAULT 0,
			CONSTRAINT party_fg FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS stats (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			created INTEGER NOT NULL,
			games_played INTEGER DEFAULT 0,
			user_registered INTEGER DEFAULT 0
		);`
	_, err := db.Exec(sql, nil)
	return err
}

func SchemaDown(db *sql.DB) error {
	sql := `
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS parties;
		DROP TABLE IF EXISTS stats;`
	_, err := db.Exec(sql, nil)
	return err
}
