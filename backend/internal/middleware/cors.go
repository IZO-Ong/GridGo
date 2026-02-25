// Package middleware provides HTTP interceptors for authentication,
// security headers, and cross-origin resource sharing.
package middleware

import (
	"net/http"
	"os"
)

// EnableCORS configures the server to allow cross-origin requests from the frontend.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Restrict origins to the FRONTEND_URL environment variable.
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		// Handle browser preflight checks
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
