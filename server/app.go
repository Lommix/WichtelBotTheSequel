package server

import (
	"database/sql"
	"errors"
	"fmt"
	"lommix/wichtelbot/server/components"
	"lommix/wichtelbot/server/store"
	"net/http"
	"time"
)

type RunState int

const (
	Debug RunState = iota
	Prod
)

const SnippetPath string = "snippets.json"

type AppState struct {
	Db        *sql.DB
	Templates *components.Templates
	Sessions  *components.CookieJar
	Snippets  *components.Snippets
	Mode      RunState
}

func (app *AppState) RegisterHandler() {
	// pages
	http.HandleFunc("/profile", app.Profile)
	http.HandleFunc("/", app.Create)
	http.HandleFunc("/login/", app.Login)
	http.HandleFunc("/join/", app.Join)

	// static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", fs)

	// htmx api
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/user", app.User)
	http.HandleFunc("/roll", app.RollDice)
	http.HandleFunc("/ping", app.PingParty)
	http.HandleFunc("/blacklist", app.GetBlacklistOptions)
}


// Game and Session Garbage Collector
func (app *AppState) CleanupRoutine() {
	for {
		time.Sleep(time.Second * 10)
		// cleaning up any left over game sessions
		expiredParties, err := store.FindExpiredParties(app.Db)
		if err != nil {
			panic(err)
		}

		if len(expiredParties) > 0 {
			fmt.Printf("cleaned %d expired parties \n", len(expiredParties))
			for _, session := range expiredParties {

				err = store.DeleteUsersInParty(app.Db, session.Id)
				if err != nil {
					panic(err)
				}

				err = session.Delete(app.Db)
				if err != nil {
					panic(err)
				}
			}
		}

		// cleaning session memeory
		count := app.Sessions.CleanupExpired()
		if count > 0 {
			fmt.Printf("cleaned %d cookies\n", count)
		}
	}
}

func (app *AppState) CurrentUserFromSession(request *http.Request) (store.User, error) {
	var user store.User
	cookie, err := request.Cookie("user")
	if err != nil {
		return user, err
	}

	for _, session := range app.Sessions.Store {
		if session.Key == cookie.Value {
			return store.FindUserWithPartyFast(app.Db,session.UserId)
		}
	}

	return user, errors.New("Not Found")
}
