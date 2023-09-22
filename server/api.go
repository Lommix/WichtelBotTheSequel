package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lommix/wichtelbot/server/components"
	"lommix/wichtelbot/server/store"
	"net/http"
)

// ----------------------------------
// Login
func Login(app *AppState, writer http.ResponseWriter, request *http.Request) {

	type LoginForm struct {
		Username string `required:"true"`
		Password string `required:"true"`
		RoomKey  string `required:"true"`
	}

	lang := components.LangFromRequest(request)
	form := &LoginForm{}
	components.FromFormData(request, form)

	if len(form.Username) == 0 || len(form.Password) == 0 {
		http.Error(writer, "invalid data", http.StatusBadRequest)
		return
	}

	user, err := store.FindUserByNameAndRoomKey(app.Db,form.Username, form.RoomKey)
	if err != nil {
		msq, _ := app.Snippets.Get("error_party_expired", lang)
		http.Error(writer, msq, http.StatusConflict)
		return
	}

	hash := sha256.Sum256([]byte(form.Password))
	pw := hex.EncodeToString(hash[:])

	if user.Password != pw {
		msq, _ := app.Snippets.Get("error_credentials", lang)
		http.Error(writer, msq, http.StatusUnauthorized)
		return
	}

	session, err := app.Sessions.CreateSession(user.Id)
	cookie := session.IntoCookie()
	http.SetCookie(writer, &cookie)

	writer.Header().Add("HX-Redirect", "/profile")
	writer.Write([]byte("ok"))
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

	writer.Header().Add("HX-Redirect", "/login/" + user.Party.Key)
}

// ----------------------------------
// Logout
func (app *AppState) User(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPut {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type updateForm struct {
		Notice    string
		ExcludeId int
	}
	form := &updateForm{}
	err := components.FromFormData(request, form)
	if err != nil {
		http.Error(writer, "invalid data", http.StatusBadRequest)
		return
	}

	var context components.TemplateContext
	context.User, err = app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "forbidden", http.StatusUnauthorized)
		return
	}


	lang := components.LangFromRequest(request)
	context.Snippets = app.Snippets.GetList(lang)
	context.User.Notice = form.Notice
	context.User.ExcludeId = int64(form.ExcludeId)

	err = context.User.Update(app.Db)
	if err != nil {
		http.Error(writer, "something went wrong", http.StatusBadRequest)
		return
	}

	err = app.Templates.Render(writer, "user", context)
	if err != nil {
		http.Error(writer, "unknown template", http.StatusBadRequest)
		return
	}
}

// ----------------------------------
// Ping Party
func (app *AppState) PingParty(writer http.ResponseWriter, request *http.Request) {
	var err error
	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "not authorized", http.StatusUnauthorized)
		return
	}

	type userState struct {
		Blacklist bool
	}

	var formData userState
	err = components.FromFormData(request, &formData)
	if err != nil {
		http.Error(writer, "invalid data", http.StatusBadRequest)
		return
	}

	if user.Party.State == store.Played || user.Party.Blacklist != formData.Blacklist {
		writer.Header().Add("HX-Refresh", "true")
	}

	writer.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(writer, "%d", len(*user.Party.Users))
}

// ----------------------------------
// Get Blacklist options live
func (app *AppState) Blacklist(writer http.ResponseWriter, request *http.Request) {
	var context components.TemplateContext
	var err error
	context.User, err = app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "not authorized", http.StatusUnauthorized)
		return
	}

	switch request.Method {
	case http.MethodGet:
		err = app.Templates.Render(writer, "blacklistOptions", context)
		if err != nil {
			http.Error(writer, "unknown template", http.StatusBadRequest)
			return
		}
	case http.MethodPost:
		type rollPostData struct {
			Blacklist bool
		}
		var formData rollPostData
		err = components.FromFormData(request, &formData)
		if err != nil {
			http.Error(writer, "invalid data", http.StatusBadRequest)
			return
		}

		context.User.Party.Blacklist = formData.Blacklist
		context.User.Party.Update(app.Db)
		writer.Header().Add("HX-Refresh", "true")
	}
}

// ----------------------------------
// Play
func (app *AppState) RollDice(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	lang := components.LangFromRequest(request)
	user, err := app.CurrentUserFromSession(request)
	if err != nil || user.Role != store.Moderator {
		msq, _ := app.Snippets.Get("error_credentials", lang)
		http.Error(writer, msq, http.StatusUnauthorized)
		return
	}

	err = user.Party.RollPartners(app.Db)
	if err != nil {
		msq, _ := app.Snippets.Get("error_roll", lang)
		http.Error(writer, msq, http.StatusExpectationFailed)
		return
	}

	store.AddGamePlayed(app.Db)
	writer.Header().Add("HX-Refresh", "true")
}
