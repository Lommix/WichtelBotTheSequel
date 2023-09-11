package server

import (
	"database/sql"
	"net/http"
)

type RunState int

const (
	Debug RunState = iota
	Prod
)

type AppState struct {
	Db       *sql.DB
	Tmpl     Templates
	Mode     RunState
	Sessions CookieJar
}

func (app *AppState) ListenAndServe(adr string) {
	// pages
	http.HandleFunc("/", app.Home)
	http.HandleFunc("/profile", app.Profile)
	// static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// api
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/user", app.User)

	println("staring server, listing on: ", adr)
	http.ListenAndServe(adr, nil)
}

// ----------------------------------
// Pages
func (app *AppState) Home(writer http.ResponseWriter, request *http.Request) {

	if app.Mode == Debug {
		app.Tmpl.Refresh()
	}

	err := app.Tmpl.store.ExecuteTemplate(writer, "home.html", nil)
	if err != nil {
		println(err.Error())
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

func (app *AppState) Profile(writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Tmpl.Refresh()
	}

	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "forbidden", http.StatusBadRequest)
		return
	}

	err = app.Tmpl.store.ExecuteTemplate(writer, "profile.html", user)
	if err != nil {
		println(err.Error())
		http.Error(writer, "forbidden", http.StatusBadRequest)
	}
}

// ----------------------------------
// Api
func (app *AppState) Logout(writer http.ResponseWriter, request *http.Request) {

	cookie := http.Cookie{
		Name:     "user",
		Value:    "",
	}

	http.SetCookie(writer, &cookie)

	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "cannot logout what is not loged in", http.StatusBadRequest)
		return
	}
	app.Sessions.DeleteSession(user.Id)
	writer.Write([]byte("logout success"))
}

func (app *AppState) User(writer http.ResponseWriter, request *http.Request) {

	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "forbidden", http.StatusBadRequest)
		return
	}

	err = app.Tmpl.store.ExecuteTemplate(writer, "user.html", user)
	if err != nil {
		http.Error(writer, "forbidden", http.StatusBadRequest)
		return
	}
}
