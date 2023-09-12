package server

import (
	"net/http"
)

// ----------------------------------
// User endpoint
func (app *AppState) User(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		err := userGet(app, writer, request)
		if err != nil {
			println(err.Error())
			http.Error(writer, "forbidden", http.StatusForbidden)
		}
	case http.MethodPut:
		err := userPut(app, writer, request)
		if err != nil {
			println(err.Error())
			http.Error(writer, "forbidden", http.StatusForbidden)
		}
	default:
		http.Error(writer, "forbidden", http.StatusMethodNotAllowed)
		return
	}
}

type updateForm struct {
	Notice    string
	Allergies string
	Exclude   string
}

func userPut(app *AppState, writer http.ResponseWriter, request *http.Request) error {
	form := &updateForm{}
	err := FromFormData(request, form)
	if err != nil {
		return err
	}


	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		return err
	}

	user.Notice = form.Notice
	user.Allergies = form.Allergies

	err = user.Update(app.Db)
	if err != nil {
		return err
	}

	err = app.Templates.Render(writer, "user", user)
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

	err = app.Templates.Render(writer, "user", user)
	if err != nil {
		return err
	}
	return nil
}
