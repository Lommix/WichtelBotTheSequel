package server

import (
	"database/sql"
	"net/http"
	"strings"
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
// Pages
func (app *AppState) Home(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if path != "/" && !strings.HasPrefix(path, "/key/") {
		return
	}

	if app.Mode == Debug {
		app.Tmpl.Load()
	}

	err := app.Tmpl.Render(writer, "home.html", nil)
	if err != nil {
		println(err.Error())
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

func (app *AppState) Profile(writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Tmpl.Load()
	}

	println("requesting profile")

	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	err = app.Tmpl.Render(writer, "profile.html", user)

	if err != nil {
		println(err.Error())
		http.Error(writer, "forbidden", http.StatusBadRequest)
	}
}

// ----------------------------------
// Api
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

	// http.Redirect(writer, request, "/login", http.StatusFound)
	writer.Header().Add("HX-Redirect","/profile")
}

func (app *AppState) User(writer http.ResponseWriter, request *http.Request) {
	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "forbidden", http.StatusBadRequest)
		return
	}

	err = app.Tmpl.Render(writer, "user", user)
	if err != nil {
		http.Error(writer, "forbidden", http.StatusBadRequest)
		return
	}

}
