package server

import (
	"crypto/sha256"
	"fmt"
	"lommix/wichtelbot/server/components"
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
	err := components.FromFormData(request, &formData)
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

	var role store.UserRole
	var party store.Party

	if len(roomKey) == 0 {
		party, err = store.CreateParty(app.Db)
		role = store.Moderator
		if err != nil {
			fmt.Println(err.Error())
			http.Error(writer, "something went wrong", http.StatusBadRequest)
			return
		}
	} else {
		party, err = store.FindPartyByKey(roomKey, app.Db)
		role = store.DefaultUser
		if err != nil {
			http.Error(writer, "invalid room", http.StatusBadRequest)
			return
		}
		if party.State == store.Played {
			http.Error(writer, "The party played without you", http.StatusBadRequest)
			return
		}
	}

	hash := sha256.Sum256([]byte(formData.Password))
	user, err := store.CreateUser(
		app.Db,
		party.Id,
		formData.Username,
		fmt.Sprintf("%x", hash),
		"",
		role,
	)

	session, err := app.Sessions.CreateSession(user.Id)
	cookie := session.IntoCookie()
	http.SetCookie(writer, &cookie)
	writer.Header().Add("HX-Redirect","/profile")
}
