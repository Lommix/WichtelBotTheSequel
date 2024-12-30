package server

import (
	// "fmt"
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

	lang := components.LangFromRequest(request)
	var err error
	context := components.TemplateContext{}
	context.Snippets = app.Snippets.GetList(lang)
	context.User, _ = app.CurrentUserFromSession(request)

	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.

	if context.IsLoggedIn() {
		http.Redirect(writer, request, "/profile", http.StatusMovedPermanently)
		return
	}

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

	lang := components.LangFromRequest(request)
	var err error
	context := components.TemplateContext{}
	context.Snippets = app.Snippets.GetList(lang)
	context.RoomKey = strings.TrimPrefix(request.URL.Path, "/join/")

	err = app.Templates.Render(writer, "join.html", context)
	if err != nil {
		http.Error(writer, "Not found", http.StatusNotFound)
	}

	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
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

		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.

		lang := components.LangFromRequest(request)
		var err error
		context := components.TemplateContext{}
		context.Snippets = app.Snippets.GetList(lang)
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

	lang := components.LangFromRequest(request)
	var err error
	context := components.TemplateContext{}
	context.Snippets = app.Snippets.GetList(lang)
	context.User, _ = app.CurrentUserFromSession(request)

	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.

	if !context.IsLoggedIn() {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
		return
	}

	err = app.Templates.Render(writer, "profile.html", context)
	if err != nil {
		http.Error(writer, "Not found", http.StatusNotFound)
	}

}
