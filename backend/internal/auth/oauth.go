package auth

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

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