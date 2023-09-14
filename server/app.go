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

type Language string

const (
	German  Language = "de"
	English Language = "en"
)

const SnippetPath string = "snippets.json"

type AppState struct {
	Db        *sql.DB
	Templates *components.Templates
	Sessions  *components.CookieJar
	Mode      RunState
}

func (app *AppState) ListenAndServe(adr string) {
	// pages
	http.HandleFunc("/profile", app.Profile)
	http.HandleFunc("/", app.Create)
	http.HandleFunc("/login/", app.Login)
	http.HandleFunc("/join/", app.Join)

	// static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// api
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/user", app.User)
	http.HandleFunc("/roll", app.RollDice)

	println("staring server, listing on: ", adr)
	http.ListenAndServe(adr, nil)
}


// Game and Session Garbage Collector
func (app *AppState) CleanupRoutine() {
	for {

		time.Sleep(time.Minute)

		// cleaning up any left over game sessions
		expiredSessions, err := store.FindExpiredParties(app.Db)
		if err != nil {
			panic(err)
		}

		if len(expiredSessions) > 0 {
			fmt.Printf("Cleaning %d expired sessions\n", len(expiredSessions))
			for _, session := range expiredSessions {
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
		app.Sessions.CleanupExpired()
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
			return store.FindUserById(session.UserId, app.Db)
		}
	}

	return user, errors.New("Not Found")
}

// ----------------------------------
// helper function

func (app *AppState) defaultContext(writer http.ResponseWriter, request *http.Request) *components.TemplateContext {
	var context components.TemplateContext
	user, err := app.CurrentUserFromSession(request)
	if err == nil {
		context.User = user
		session, err := store.FindPartyByID(user.PartyId, app.Db)
		if err == nil {
			context.User.GameSession = &session
		}
	}

	// todo cache this, add lang select
	snippets, err := components.LoadSnippets(string(German), SnippetPath)
	context.Snippets = snippets

	return &context
}
