package store

import "database/sql"

func SchemaUp(db *sql.DB) error {
	sql := `
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			created INTEGER NOT NULL,
			state INTEGER NOT NULL,
			key STRING NOT NULL,
			rule_set INTEGER DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			session_id INTEGER NOT NULL,
			created INTEGER NOT NULL,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			partner_id INTEGER DEFAULT 0,
			exclude_id TEXT DEFAULT 0,
			notice TEXT DEFAULT NULL,
			role INTEGER DEFAULT 0,
			FOREIGN KEY (session_id) REFERENCES sessions(id)
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
		DROP TABLE IF EXISTS sessions;
		DROP TABLE IF EXISTS stats;`
	_, err := db.Exec(sql, nil)
	return err
}
