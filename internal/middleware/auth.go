package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var TokenSecret = []byte("Super Secret Key")

func Auth(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			logger.Error(reqIDString + "Missing token header")
			next.ServeHTTP(w, r)

			return
		}

		tokenParts := strings.Split(authToken, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Error(reqIDString+"Invalid token header format: ", authToken)
			http.Error(w, "Invalid token header format", http.StatusUnauthorized)

			return
		}

		authToken = tokenParts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
			return TokenSecret, nil
		})

		if err != nil || !token.Valid {
			logger.Errorf(reqIDString+"Auth err: %s, is token valid: %t", err, token.Valid)
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.Error(reqIDString + "User ID not found in token")
			http.Error(w, "User ID not found in token", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
