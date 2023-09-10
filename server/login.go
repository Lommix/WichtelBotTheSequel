package server

import (
	"crypto/sha256"
	"encoding/hex"
	"lommix/wichtelbot/server/store"
	"net/http"
)

type LoginForm struct {
	Username string
	Password string
	RoomKey  string
}

func (app *AppState) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}

	form := &LoginForm{}
	ParseFormInto(request, form)

	if len(form.Username) == 0 || len(form.Password) == 0 {
		http.Error(writer, "invalid data", http.StatusBadRequest)
		return
	}

	user, err := store.FindUserByNameAndRoomKey(form.Username, form.RoomKey, app.Db)
	if err != nil {
		http.Error(writer, "invalid user", http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256([]byte(form.Password))
	pw := hex.EncodeToString(hash[:])

	if user.Password != pw {
		http.Error(writer, "wrong password", http.StatusBadRequest)
		return
	}

	cookie, err := app.Sessions.CreateSession(user.Id)

	writer.Header().Add("Set-Cookie", cookie.IntoCookie())
	writer.Write([]byte("ok"))
}
