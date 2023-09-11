package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"lommix/wichtelbot/server"
	"lommix/wichtelbot/server/store"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		return
	}
	command := os.Args[1]


	db, err := sql.Open("sqlite3", "wichtel.db")
	if err != nil {
		fmt.Println("Failed to open DB Connection: ", err)
		return
	}

	switch command {
	case "init":
		store.SchemaDown(db)
		store.SchemaUp(db)
	case "dev":
		tmpl := server.Templates{
			Dir: "./templates",
		}

		err = tmpl.Load()

		if err != nil {
			fmt.Println("Failed to load Templates: ", err)
			return
		}

		app := server.AppState{
			Db:   db,
			Tmpl: tmpl,
			Mode: server.Debug,
			Sessions: server.CookieJar{},
		}

		println("starting cleaner")
		go app.CleanupRoutine()

		app.ListenAndServe(":3000")

	case "prod":
		fmt.Println("Not implemented yet")
	default:
		fmt.Println("Invalid command\nOptions are\n'init' 'dev' 'prod' ")
	}
}

func loadTemplates() (*template.Template, error) {
	tmpl, err := template.ParseFiles(
		"./templates/components/layout.html",
		"./templates/home.html",
		"./templates/profile.html",
	)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
