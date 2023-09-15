package server

import (
	"lommix/wichtelbot/server/components"
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
		app.Snippets.Load()
	}

	var err error
	context := components.TemplateContext{}
	context.Snippets = app.Snippets.GetList(components.German)
	context.User, _ = app.CurrentUserFromSession(request)

	err = app.Templates.Render(writer, "create.html", context)
	if err != nil {
		http.Error(writer, "Not found", http.StatusNotFound)
	}
}

// ----------------------------------
// join page
func (app *AppState) Join(writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Templates.Load()
		app.Snippets.Load()
	}

	var err error
	context := components.TemplateContext{}
	context.Snippets = app.Snippets.GetList(components.German)

	err = app.Templates.Render(writer, "join.html", context)
	if err != nil {
		http.Error(writer, "Not found", http.StatusNotFound)
	}
}

// ----------------------------------
// login page
func (app *AppState) Login(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:

		if app.Mode == Debug {
			app.Templates.Load()
			app.Snippets.Load()
		}


		var err error
		context := components.TemplateContext{}
		context.Snippets = app.Snippets.GetList(components.German)
		context.User, _ = app.CurrentUserFromSession(request)

		if context.IsLoggedIn() {
			http.Redirect(writer, request, "/profile", http.StatusMovedPermanently)
			return
		}

		err = app.Templates.Render(writer, "login.html", context)
		if err != nil {
			http.Error(writer, "Not found", http.StatusNotFound)
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

	var err error
	context := components.TemplateContext{}
	context.Snippets = app.Snippets.GetList(components.German)
	context.User, _ = app.CurrentUserFromSession(request)

	if !context.IsLoggedIn() {
		http.Redirect(writer, request, "/login", http.StatusUnauthorized)
		return
	}

	err = app.Templates.Render(writer, "profile.html", context)
	if err != nil {
		http.Error(writer, "Not found", http.StatusNotFound)
	}
}
