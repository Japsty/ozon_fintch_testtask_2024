package services

import (
	"Ozon_testtask/internal/model"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrInvalidToken  = errors.New("token not valid")
	TokenSecret      = []byte("Super Secret Key")
)

type SessionService struct {
	sessionRepo model.SessionRepo
}

func NewSessionService(sessionRepo model.SessionRepo) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (ss *SessionService) GetTokenData(r *http.Request) (*model.Session, error) {
	inToken := r.Header.Get("Authorization")
	if inToken == "" {
		return nil, ErrTokenNotFound
	}

	tokenParts := strings.Split(inToken, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, ErrInvalidToken
	}

	token := tokenParts[1]

	session, err := ss.sessionRepo.GetSession(token)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return session, nil
}

func (ss *SessionService) SetUpJWT(user model.User) ([]byte, error) {
	timeToExpire := time.Now().Add(48 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": user.Username,
			"id":       user.ID,
		},
		"exp": timeToExpire,
	})

	tokenString, err := token.SignedString(TokenSecret)
	if err != nil {
		return nil, err
	}

	err = ss.sessionRepo.Create(timeToExpire, tokenString, user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
