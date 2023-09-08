package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type AppState struct {
	db *sql.DB
}

func main() {

	db, err := sql.Open("sqlite3", "wichtel.db")

	if err != nil {
		fmt.Println("Err: ", err)
		return
	}

	test, _ := db.Exec("SHOW TABLES;")

	fmt.Print(test)
}

func (app *AppState) handle() {}
