package server

import (
	"crypto/sha256"
	"encoding/hex"
	"lommix/wichtelbot/server/store"
	"net/http"
)

type LoginForm struct {
	Username string `required:"true"`
	Password string `required:"true"`
	RoomKey  string `required:"true"`
}

// ----------------------------------
// login controller
func (app *AppState) Login(writer http.ResponseWriter, request *http.Request) {
	println("requesting login")

	switch request.Method {
	case http.MethodGet:
		loginGet(app, writer, request)
	case http.MethodPost:
		loginPost(app, writer, request)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ----------------------------------
// login get
func loginGet(app *AppState, writer http.ResponseWriter, request *http.Request) {
	if app.Mode == Debug {
		app.Templates.Load()
	}

	err := app.Templates.Render(writer, "login.html", app.defaultContext(writer, request))
	if err != nil {
		println(err.Error())
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}
}

// ----------------------------------
// login post
func loginPost(app *AppState, writer http.ResponseWriter, request *http.Request) {
	form := &LoginForm{}
	FromFormData(request, form)

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
