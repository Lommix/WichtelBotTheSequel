package main

import (
	"database/sql"
	"fmt"
	"lommix/wichtelbot/server"
	"lommix/wichtelbot/server/components"
	"lommix/wichtelbot/server/store"
	"net/http"
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
	// -------------------------------------
	// init database
	// -------------------------------------
	case "init":
		store.SchemaDown(db)
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

		app := server.AppState{
			Db:        db,
			Templates: tmpl,
			Mode:      server.Debug,
			Snippets:  snippets,
			Sessions:  &components.CookieJar{},
		}
		println("starting cleaner")
		go app.CleanupRoutine()

		println("starting listeneing on 3000")
		app.RegisterHandler()
		http.ListenAndServe(":3000", nil)
	// -------------------------------------
	// start production server in tls mode
	// -------------------------------------
	case "prod":
		components.LoadEnv()
		cert := os.Getenv("SSL_CERT")
		key := os.Getenv("SSL_KEY")
		http_port := os.Getenv("HTTP_PORT")
		https_port := os.Getenv("HTTPS_PORT")

		if cert == "" || key == "" || http_port == "" || https_port == "" {
			fmt.Println("Please provide SSL_CERT, SSL_KEY, HTTP_PORT, HTTPS_PORT")
			return
		}

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

		app := server.AppState{
			Db:        db,
			Templates: tmpl,
			Mode:      server.Debug,
			Snippets:  snippets,
			Sessions:  &components.CookieJar{},
		}

		println("starting cleaner")
		go app.CleanupRoutine()


		println("starting http")
		go http.ListenAndServe(":"+http_port, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			http.Redirect(writer, request, "https://"+request.Host, http.StatusMovedPermanently)
		}))

		println("starting https")
		app.RegisterHandler()
		err := http.ListenAndServeTLS(":"+https_port, "localhost.crt", "localhost.key", nil)
		if err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("Invalid command\nOptions are\n'init' 'dev' 'prod' ")
	}
}

func redirectHttp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}
