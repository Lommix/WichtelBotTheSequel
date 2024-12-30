package components

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

type UserRole int

const (
	Guest UserRole = iota
	User
	Moderator
)
const CookieExpirationaTime time.Duration = time.Hour * 24 * 30

type Session struct {
	UserId  int64
	Created int64
	Key     string
}

type CookieJar struct {
	Store []Session
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
	session.Key = hex.EncodeToString(hash[:])
	session.Created = time.Now().Unix()

	jar.Store = append(jar.Store, session)

	return session, nil
}

func (session *Session) IntoCookie() http.Cookie {
	return http.Cookie{
		Name:     "user",
		Value:    session.Key,
		Path:     "/",
		Expires:  time.Now().Add(CookieExpirationaTime),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
}

func (app *CookieJar) DeleteSession(userId int64) error {
	for i := 0; i < len(app.Store); i++ {
		if app.Store[i].UserId == userId {
			app.Store = append(app.Store[:i], app.Store[i+1:]...)
			i--
		}
	}
	return nil
}

func (app *CookieJar) CleanupExpired() int {
	count := 0
	timeOffset := time.Now().Add(-CookieExpirationaTime).Unix()

	for i := 0; i < len(app.Store); i++ {
		if app.Store[i].Created < timeOffset {
			app.Store = append(app.Store[:i], app.Store[i+1:]...)
			count++
			i--
		}
	}
	return count
}
