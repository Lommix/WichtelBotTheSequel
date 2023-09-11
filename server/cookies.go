package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserRole int

const (
	Guest UserRole = iota
	User
	Moderator
)
const CookieExpirationaTime time.Duration = time.Hour * 24 * 4

type Session struct {
	UserId  int64
	Created time.Time
	key     string
}

type CookieJar struct {
	store []Session
}

func (jar *CookieJar) FindByKey(key []byte) (Session, error) {
	var session Session
	return session, nil
}

func (jar *CookieJar) CreateSession(userId int64) (Session, error) {
	var session Session
	session.UserId = userId

	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return session, err
	}

	hash := sha256.Sum256(randomBytes)
	session.key = hex.EncodeToString(hash[:])
	session.Created = time.Now()

	return session, nil
}

func (session *Session) IntoCookie() string {
	var out string
	out = fmt.Sprintf("user=%s", session.key)
	return out
}

func (app *AppState) UserIdFromRequest(request *http.Request) (int64, error) {
	cookie, err := request.Cookie("user")
	if err != nil {
		return 0, err
	}

	for _, session := range app.Sessions.store {
		if session.key == cookie.Value {
			return session.UserId, nil
		}
	}

	return 0, errors.New("Not Found")
}

func (app *CookieJar) CleanupExpired() {
	for i := 0; i < len(app.store); i++ {
		if app.store[i].Created.Unix() > time.Now().Add(-CookieExpirationaTime).Unix() {
			app.store = append(app.store[:i], app.store[i+1:]...)
			i--
		}
	}
}
