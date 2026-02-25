// Package middleware provides HTTP interceptors for authentication,
// security headers, and cross-origin resource sharing.
package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type used to prevent key collisions in context.Context.
type contextKey string

// UserIDKey is the unique key used to store/retrieve the user's ID from the request context.
const UserIDKey contextKey = "user_id"

// GetJWTKey retrieves the secret string used for signing and verifying JSON Web Tokens.
func GetJWTKey() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// GetUserID is a helper function that extracts the authenticated user's ID 
// from the request context. Returns an empty string if the user is not authenticated.
func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDKey).(string)
	return id
}

// OptionalAuth attempts to authenticate the user via a Bearer token.
// If the token is valid, the UserID is added to the context.
// If the token is missing or invalid, the request still proceeds (unauthenticated).
func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		// If no Bearer token is provided, just move to the next handler.
		if !strings.HasPrefix(authHeader, "Bearer ") {
			next(w, r)
			return
		}

		tokenString := authHeader[7:]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return GetJWTKey(), nil
		})

		// If valid, inject ID into context.
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if id, ok := claims["user_id"].(string); ok {
					ctx := context.WithValue(r.Context(), UserIDKey, id)
					next(w, r.WithContext(ctx))
					return
				}
			}
		}
		
		// Proceed even if token was invalid (optional auth).
		next(w, r)
	}
}

// RequireAuth enforces authentication.
// If a Bearer token is not present, it returns a 401 Unauthorized error.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "AUTHENTICATION_REQUIRED", http.StatusUnauthorized)
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

		http.Error(w, "INVALID_OR_EXPIRED_TOKEN", http.StatusUnauthorized)
	}
}