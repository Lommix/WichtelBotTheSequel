package server

import (
	"crypto/sha256"
	"fmt"
	"lommix/wichtelbot/server/store"
	"net/http"
)

type RegisterForm struct {
	Username string
	Password string
	Retry    string
	roomKey  string
}

func (app *AppState) Register(writer http.ResponseWriter, request *http.Request) {
	var formData RegisterForm
	err := ParseFormInto(request, &formData)
	if err != nil {
		http.Error(writer, "Invalid post", http.StatusBadRequest)
		return
	}

	if formData.Password != formData.Retry {
		http.Error(writer, "Passwords not matching", http.StatusBadRequest)
		return
	}


	roomKey := func() string {
		if len(formData.roomKey) > 0 {
			return formData.roomKey
		}
		queryParams := request.URL.Query()
		if len(queryParams.Get("roomKey")) > 0 {
			return queryParams.Get("roomKey")
		}
		return ""
	}()

	var session store.GameSession
	if len(roomKey) == 0 {
		session, err = store.CreateSession(app.Db)
		if err != nil {
			http.Error(writer, "invalid post", http.StatusBadRequest)
			return
		}
	} else {
		session, err = store.FindSessionByKey(formData.roomKey, app.Db)
		if err != nil {
			http.Error(writer, "invalid post", http.StatusBadRequest)
			return
		}
	}

	hash := sha256.Sum256([]byte(formData.Password))
	user, err := store.CreateUser(
		app.Db,
		session.Id,
		formData.Username,
		fmt.Sprintf("%x", hash),
		"",
		"",
		store.Moderator,
	)

	cookie, err := app.Sessions.CreateSession(user.Id)

	writer.Header().Add("Set-Cookie", cookie.IntoCookie())
	writer.Write([]byte("ok"))
}
