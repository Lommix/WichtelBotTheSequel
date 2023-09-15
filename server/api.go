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

	form := &LoginForm{}
	components.FromFormData(request, form)

	if len(form.Username) == 0 || len(form.Password) == 0 {
		http.Error(writer, "invalid data", http.StatusBadRequest)
		return
	}

	user, err := store.FindUserByNameAndRoomKey(app.Db,form.Username, form.RoomKey)
	if err != nil {
		msq, _ := app.Snippets.Get("error_party_expired", components.German)
		http.Error(writer, msq, http.StatusConflict)
		return
	}

	hash := sha256.Sum256([]byte(form.Password))
	pw := hex.EncodeToString(hash[:])

	if user.Password != pw {
		msq, _ := app.Snippets.Get("error_credentials", components.German)
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
	switch request.Method {
	case http.MethodGet:
		err := UserGet(app, writer, request)
		if err != nil {
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
// ----------------------------------
// Ping Party
func (app *AppState) PingParty(writer http.ResponseWriter, request *http.Request) {

	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "not authorized", http.StatusUnauthorized)
		return
	}

	if user.Party.State != store.Played {
		http.Error(writer, "not played yet", http.StatusNoContent)
		return
	}

	writer.Header().Add("HX-Refresh", "true")
}

// ----------------------------------
// Play
func (app *AppState) RollDice(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type rollPostData struct {
		Rule string
	}

	var formData rollPostData
	err := components.FromFormData(request, &formData)

	if err != nil {
		http.Error(writer, "invalid data", http.StatusBadRequest)
		return
	}

	withBlacklist := func()bool{
		if formData.Rule == "blacklist" {
			return true
		}
		return false
	}()

	fmt.Println("Blackilist:", withBlacklist)

	user, err := app.CurrentUserFromSession(request)
	if err != nil || user.Role != store.Moderator {
		msq, _ := app.Snippets.Get("error_credentials", components.German)
		http.Error(writer, msq, http.StatusUnauthorized)
		return
	}

	err = user.Party.RollPartners(app.Db, withBlacklist)

	if err != nil {
		msq, _ := app.Snippets.Get("error_roll", components.German)
		http.Error(writer, msq, http.StatusExpectationFailed)
		return
	}

	writer.Header().Add("HX-Refresh", "true")
}
