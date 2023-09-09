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
	Db   *sql.DB
	Tmpl Templates
	Mode RunState
}

func (app *AppState) ListenAndServe(adr string) {
	//pages
	http.HandleFunc("/", app.Home)
	http.HandleFunc("/profile", app.Profile)

	//static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//api

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

	err := app.Tmpl.store.ExecuteTemplate(writer, "profile.html", nil)
	if err != nil {
		println(err.Error())
		http.Error(writer, "forbidden", http.StatusBadRequest)
	}
}

// ----------------------------------
// Static
func (app *AppState) Static(writer http.ResponseWriter, request *http.Request) {

}

// ----------------------------------
// Api

func (app *AppState) Logout(writer http.ResponseWriter, request *http.Request) {
}

func (app *AppState) User(writer http.ResponseWriter, request *http.Request) {
}
