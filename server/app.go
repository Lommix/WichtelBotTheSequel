package server

import (
	"database/sql"
	"fmt"
	"lommix/wichtelbot/server/store"
	"net/http"
	"strings"
)

type RunState int

const (
	Debug RunState = iota
	Prod
)

type AppState struct {
	Db        *sql.DB
	Templates Templates
	Mode      RunState
	Sessions  CookieJar
}

func (app *AppState) ListenAndServe(adr string) {
	// pages
	http.HandleFunc("/profile", app.Profile)
	http.HandleFunc("/", app.Home)
	http.HandleFunc("/login/", app.Login)

	// static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// api
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/user", app.User)

	println("staring server, listing on: ", adr)
	http.ListenAndServe(adr, nil)
}

// ----------------------------------
// home page
func (app *AppState) Home(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if path != "/" && !strings.HasPrefix(path, "/key/") {
		return
	}

	if app.Mode == Debug {
		app.Templates.Load()
	}


	err := app.Templates.Render(writer, "home.html", app.defaultContext(writer, request))
	if err != nil {
		println(err.Error())
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

// ----------------------------------
// profile page
func (app *AppState) Profile(writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Templates.Load()
	}

	println("requesting profile")

	context := app.defaultContext(writer, request)
	fmt.Print(context.User)
	if context.User.Id == 0 {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	err := app.Templates.Render(writer, "profile.html", context)

	if err != nil {
		println(err.Error())
		http.Error(writer, "forbidden", http.StatusBadRequest)
	}
}

// ----------------------------------
// Logout
func (app *AppState) Logout(writer http.ResponseWriter, request *http.Request) {
	cookie := http.Cookie{
		Name:  "user",
		Value: "",
	}

	http.SetCookie(writer, &cookie)

	user, err := app.CurrentUserFromSession(request)
	if err == nil {
		app.Sessions.DeleteSession(user.Id)
	}

	writer.Header().Add("HX-Redirect", "/login")
}


func (app *AppState) defaultContext(writer http.ResponseWriter, request *http.Request) *TemplateContext {
	var context TemplateContext
	user, err := app.CurrentUserFromSession(request)
	if err == nil {
		context.User = user
		session,err := store.FindSessionByID(user.Session_id,app.Db)
		if err == nil {
			context.User.GameSession = &session
		}
	}

	return &context
}
