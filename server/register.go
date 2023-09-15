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
		http.Error(writer, "Invalid post", http.StatusUnauthorized)
		return
	}

	if formData.Password != formData.Retry {
		msq, _ := app.Snippets.Get("error_retry", components.German)
		http.Error(writer, msq, http.StatusUnauthorized)
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
			http.Error(writer, "Server busy", http.StatusBadRequest)
			return
		}
	} else {
		party, err = store.FindPartyByKey(roomKey, app.Db)
		role = store.DefaultUser
		msq, _ := app.Snippets.Get("error_party_expired", components.German)
		if err != nil {
			http.Error(writer, msq, http.StatusConflict)
			return
		}
		if party.State == store.Played {
			http.Error(writer, msq, http.StatusConflict)
			return
		}
	}

	for _, u := range *party.Users {
		if u.Name == formData.Username {
			msq, _ := app.Snippets.Get("error_name_taken", components.German)
			http.Error(writer, msq, http.StatusConflict)
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
	writer.Header().Add("HX-Redirect", "/profile")
}
