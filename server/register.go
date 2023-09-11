package server

import (
	"crypto/sha256"
	"fmt"
	"lommix/wichtelbot/server/store"
	"net/http"
)

type RegisterForm struct {
	Username string `required:"true"`
	Password string `required:"true"`
	Retry    string `required:"true"`
	RoomKey  string
}

func (app *AppState) Register(writer http.ResponseWriter, request *http.Request) {
	var formData RegisterForm
	err := FromFormData(request, &formData)
	if err != nil {
		println(err.Error())
		http.Error(writer, "Invalid post", http.StatusBadRequest)
		return
	}

	if formData.Password != formData.Retry {
		writer.Header().Add("HX-Target", "#create-error")
		writer.Header().Add("HX-Swap", "innerHTML")
		http.Error(writer, "Passwords not matching", http.StatusBadRequest)
		return
	}


	roomKey := func() string {
		if len(formData.RoomKey) > 0 {
			return formData.RoomKey
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
			fmt.Println(err.Error())
			http.Error(writer, "something wen wrong", http.StatusBadRequest)
			return
		}
	} else {
		session, err = store.FindSessionByKey(roomKey, app.Db)
		if err != nil {
			http.Error(writer, "invalid room", http.StatusBadRequest)
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
	writer.Header().Add("HX-Redirect","/profile")
}
