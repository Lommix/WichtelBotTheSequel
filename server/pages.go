package server

import (
	"net/http"
	"strings"
)

// ----------------------------------
// create page
func (app *AppState) Create(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if path != "/" && !strings.HasPrefix(path, "/key/") {
		return
	}

	if app.Mode == Debug {
		app.Templates.Load()
	}

	err := app.Templates.Render(writer, "create.html", app.defaultContext(request))
	if err != nil {
		println(err.Error())
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

// ----------------------------------
// join page
func (app *AppState) Join(writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Templates.Load()
	}

	err := app.Templates.Render(writer, "join.html", app.defaultContext(request))
	if err != nil {
		println(err.Error())
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

// ----------------------------------
// login page
func (app *AppState) Login(writer http.ResponseWriter, request *http.Request) {
	println("requesting login")

	switch request.Method {
	case http.MethodGet:
		if app.Mode == Debug {
			app.Templates.Load()
		}

		err := app.Templates.Render(writer, "login.html", app.defaultContext(request))
		if err != nil {
			println(err.Error())
			http.Error(writer, "Bad Request", http.StatusBadRequest)
		}
	case http.MethodPost:
		Login(app, writer, request)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ----------------------------------
// profile page
func (app *AppState) Profile(writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Templates.Load()
	}

	context := app.defaultContext(request)
	if !context.IsLoggedIn() {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	err := app.Templates.Render(writer, "profile.html", context)
	if err != nil {
		println(err.Error())
		http.Error(writer, "forbidden", http.StatusBadRequest)
	}
}
