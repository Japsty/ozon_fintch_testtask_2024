package model

import (
	"context"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash []byte
	ID           string `json:"id"`
}

type DBUser struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	ID           string `json:"id"`
}

type Session struct {
	UserID string
	Token  string
	Exp    int64
}

type UserRepo interface {
	CreateUser(context.Context, string, string, []byte) (User, error)
	GetUser(context.Context, string, string) (User, error)
}

type UserService interface {
	AddUser(context.Context, string, string) (User, error)
	GetUser(context.Context, string, string) (User, error)
}
