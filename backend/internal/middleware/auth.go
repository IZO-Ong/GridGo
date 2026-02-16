package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const UserIDKey contextKey = "user_id"

func GetJWTKey() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDKey).(string)
	return id
}

func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			next(w, r)
			return
		}

		tokenString := authHeader[7:]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return GetJWTKey(), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if id, ok := claims["user_id"].(string); ok {
					ctx := context.WithValue(r.Context(), UserIDKey, id)
					next(w, r.WithContext(ctx))
					return
				}
			}
		}
		next(w, r)
	}
}