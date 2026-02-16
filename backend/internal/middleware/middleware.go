package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 8 { return "" }

	tokenString := authHeader[7:]
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})

	if token != nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			return claims["user_id"].(string)
		}
	}
	return ""
}