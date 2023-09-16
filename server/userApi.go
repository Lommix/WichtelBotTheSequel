package server

import (
	"lommix/wichtelbot/server/components"
	"net/http"
)

// ----------------------------------
// User endpoint
func userPut(app *AppState, writer http.ResponseWriter, request *http.Request) error {
	type updateForm struct {
		Notice    string
		ExcludeId int
	}
	form := &updateForm{}
	err := components.FromFormData(request, form)
	if err != nil {
		return err
	}

	var context components.TemplateContext
	context.User, err = app.CurrentUserFromSession(request)
	if err != nil {
		return err
	}

	lang := components.LangFromRequest(request)
	context.Snippets = app.Snippets.GetList(lang)
	context.User.Notice = form.Notice
	context.User.ExcludeId = int64(form.ExcludeId)

	err = context.User.Update(app.Db)
	if err != nil {
		return err
	}

	err = app.Templates.Render(writer, "user", context)
	if err != nil {
		return err
	}

	return nil
}

func userGet(app *AppState, writer http.ResponseWriter, request *http.Request) error {
	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		return err
	}
	return app.Templates.Render(writer, "user", user)
}
