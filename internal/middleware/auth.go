package middleware

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func Auth(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			logger.Error(reqIDString + "Missing token header")
			http.Error(w, "Missing token header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authToken, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Error(reqIDString+"Invalid token header format: ", authToken)
			http.Error(w, "Invalid token header format", http.StatusUnauthorized)
			return
		}

		authToken = tokenParts[1]

		//claims := jwt.MapClaims{}
		//token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		//	//return services.TokenSecret, nil
		//})

		//if err != nil || !token.Valid {
		//	logger.Errorf(reqIDString+"Auth err: %s, is token valid: %b", err, token.Valid)
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}

		//r = r.WithContext(context.WithValue(r.Context(), "token", token))
		//next.ServeHTTP(w, r)
	})
}
