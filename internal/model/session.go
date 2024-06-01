package model

import "net/http"

type Session struct {
	ID       int
	Username string
	UserID   string
	Token    string
	Exp      int64
}

type SessionRepo interface {
	Create(int64, string, string, string) error
	GetSession(string) (*Session, error)
	Delete(string) error
}

type SessionService interface {
	GetTokenData(*http.Request) (*Session, error)
	SetUpJWT(User) ([]byte, error)
}
