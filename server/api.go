package server

import (
	"crypto/sha256"
	"encoding/hex"
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

	user, err := store.FindUserByNameAndRoomKey(form.Username, form.RoomKey, app.Db)
	if err != nil {
		println(err.Error())
		http.Error(writer, "invalid credentials", http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256([]byte(form.Password))
	pw := hex.EncodeToString(hash[:])

	if user.Password != pw {
		http.Error(writer, "invalid credentials", http.StatusBadRequest)
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

	writer.Header().Add("HX-Redirect", "/login/" + user.GameSession.Key)
}

// ----------------------------------
// Ping Party
func (app *AppState) PingParty(writer http.ResponseWriter, request *http.Request) {
	user, err := app.CurrentUserFromSession(request)
	if err != nil {
		http.Error(writer, "not authorized", http.StatusForbidden)
		return
	}

	if user.GameSession.State != store.Played {
		http.Error(writer, "not played yet", http.StatusNoContent)
		return
	}

	writer.Header().Add("HX-Refresh", "true")
}

// ----------------------------------
// Play
func (app *AppState) RollDice(writer http.ResponseWriter, request *http.Request) {
	user, err := app.CurrentUserFromSession(request)
	if err != nil || user.Role != store.Moderator {
		http.Error(writer, "not authorized", http.StatusForbidden)
		return
	}

	println("roll dice")

	err = user.GameSession.RollPartners(app.Db)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadGateway)
		return
	}

	writer.Header().Add("HX-Refresh", "true")
}
