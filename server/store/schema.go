package store

import "database/sql"

func SchemaUp(db *sql.DB){
	sql := `
		DROP TABLE IF EXISTS sessions;
		CREATE TABLE sessions (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			created INTEGER NOT NULL,
			state INTEGER NOT NULL,
			key STRING NOT NULL,
			rules STRING DEFAULT NULL
		);

		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			session_id INTEGER NOT NULL,
			created INTEGER NOT NULL,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			partner_id INTEGER DEFAULT NULL,
			exclude_id TEXT DEFAULT NULL,
			notice TEXT DEFAULT NULL,
			allergies TEXT DEFAULT NULL,
			role TEXT DEFAULT 'normal'
			FOREIGN KEY (session_id) REFERENCES sessions(id)
		);

		DROP TABLE IF EXISTS stats;
		CREATE TABLE stats (
			id INTEGER PRIMARY KEY UNIQUE NOT NULL,
			created INTEGER NOT NULL,
			games_played INTEGER DEFAULT NULL,
			user_registered INTEGER DEFAULT NULL
		);`
	db.Exec(sql,nil)
}
