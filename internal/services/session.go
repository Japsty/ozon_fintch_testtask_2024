package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrInvalidToken  = errors.New("token not valid")
	ErrEmptyPayload  = errors.New("empty payload")
	TokenSecret      = []byte("Super Secret Key")
)

type SessionService struct {
}

func NewSessionService() *SessionService {
	return &SessionService{}
}

func (ss *SessionService) GetTokenData(r *http.Request) (string, error) {
	inToken := r.Header.Get("Authorization")
	if inToken == "" {
		return "", ErrTokenNotFound
	}

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return TokenSecret, nil
	}

	tokenParts := strings.Split(inToken, " ")
	token, err := jwt.Parse(tokenParts[1], hashSecretGetter)
	if err != nil {
		return "", err
	} else if !token.Valid {
		return "", ErrInvalidToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrEmptyPayload
	}

	userID, ok := payload["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found")
	}
	return userID, nil
}
