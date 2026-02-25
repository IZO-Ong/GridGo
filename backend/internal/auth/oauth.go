// Package auth handles the configuration and initialization of external
// authentication providers (OAuth2).
package auth

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

// NewAuth initializes the Goth authentication providers.
// It reads credentials from environment variables and defines the
// scopes (permissions) requested from the provider.
func NewAuth() {
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("BACKEND_URL")+"/api/auth/google/callback",
			"email",
			"profile",
		),
	)
}