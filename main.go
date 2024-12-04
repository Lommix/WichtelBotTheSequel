package main

import (
	"database/sql"
	"fmt"
	"lommix/wichtelbot/server"
	"lommix/wichtelbot/server/components"
	"lommix/wichtelbot/server/store"
	"net/http"
	"os"
	"strconv"

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
	// -------------------------------------
	// database commands
	// -------------------------------------
	case "init":
		store.SchemaUp(db)
	// -------------------------------------
	// start dev server
	// -------------------------------------
	case "dev":
		tmpl := &components.Templates{}
		err = tmpl.Load()
		if err != nil {
			fmt.Println("Failed to load Templates: ", err)
			return
		}

		snippets := &components.Snippets{}
		err = snippets.Load()
		if err != nil {
			fmt.Println("Failed to load Snippets: ", err)
			return
		}

		components.LoadEnv()
		settings, err := components.LoadSettingsFromEnv()
		if err != nil {
			fmt.Println("Failed to load settings: ", err)
		}

		app := server.AppState{
			Db:        db,
			Templates: tmpl,
			Mode:      server.Debug,
			Snippets:  snippets,
			Settings:  settings,
			Sessions:  &components.CookieJar{},
		}
		println("starting cleaner")
		go app.CleanupRoutine()

		println("starting listeneing on 3000")
		app.RegisterHandler()
		http.ListenAndServe(":3000", nil)
	// -------------------------------------
	// start production server
	// -------------------------------------
	case "prod":
		tmpl := &components.Templates{}
		err = tmpl.Load()
		if err != nil {
			fmt.Println("Failed to load Templates: ", err)
			return
		}

		snippets := &components.Snippets{}
		err = snippets.Load()
		if err != nil {
			fmt.Println("Failed to load Snippets: ", err)
			return
		}

		components.LoadEnv()
		settings, err := components.LoadSettingsFromEnv()
		if err != nil {
			fmt.Println("Failed to load settings: ", err)
		}

		app := server.AppState{
			Db:        db,
			Templates: tmpl,
			Snippets:  snippets,
			Sessions:  &components.CookieJar{},
			Settings:  settings,
			Mode:      server.Prod,
		}

		println("starting cleaner")
		go app.CleanupRoutine()

		if settings.Https != nil {
		println("starting http redirect")
			go http.ListenAndServe(":"+strconv.Itoa(settings.Http.Port), http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, "https://"+request.Host, http.StatusMovedPermanently)
		}))

		println("starting https")
		app.RegisterHandler()
			err := http.ListenAndServeTLS(":"+strconv.Itoa(settings.Https.Port), settings.Https.SslCertPath, settings.Https.SslKeyPath, nil)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			println("starting http")
			app.RegisterHandler()
			err := http.ListenAndServe(":"+strconv.Itoa(settings.Http.Port), nil)
		if err != nil {
			fmt.Println(err)
		}
		}
	default:
		fmt.Println("Invalid command\nOptions are\n'init' 'dev' 'prod' ")
	}
}
