package server

type UserRole int
const (
    Guest UserRole = iota
    User
    Moderator
)

type Session struct {
	UserId int;
	key []byte;
}

type CookieJar struct{
	store []Session
}
