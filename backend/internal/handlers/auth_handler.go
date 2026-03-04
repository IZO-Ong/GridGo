// Package handlers contains the HTTP logic for the GridGo API.
// This file specifically manages user lifecycle: Authentication, Registration, and OAuth.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/IZO-Ong/gridgo/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

// HandleLogin authenticates a user using either their username or email.
// On success, it returns a JWT valid for 12 hours.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Identifier string `json:"username"` // Can be username or email
		Password   string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "INVALID_PAYLOAD", http.StatusBadRequest)
		return
	}

	if creds.Identifier == "" || creds.Password == "" {
		http.Error(w, "CREDENTIALS_REQUIRED", http.StatusBadRequest)
		return
	}

	var user models.User
	// Check against both username and email fields
	if err := db.DB.Where("username = ? OR email = ?", creds.Identifier, creds.Identifier).First(&user).Error; err != nil {
		http.Error(w, "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}

	// Generate JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	
	tokenString, _ := token.SignedString(middleware.GetJWTKey())
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "username": user.Username})
}

// HandleOAuthLogin initiates the Google OAuth flow via Gothic.
func HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Set("provider", "google")
	r.URL.RawQuery = q.Encode()
	
	gothic.BeginAuthHandler(w, r)
}

// HandleOAuthCallback handles the return from Google. It creates a user 
// record if it's the first time they log in via Google, then redirects back to the frontend
func HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
    user, err := gothic.CompleteUserAuth(w, r)
    if err != nil { return }

    var dbUser models.User
    // FirstOrCreate to atomically handle the check or create logic
    err = db.DB.Where(models.User{Email: user.Email}).
        Attrs(models.User{
            Username:     user.NickName, 
            PasswordHash: "OAUTH_ACCOUNT",
        }).
        FirstOrCreate(&dbUser).Error

    if err != nil {
        http.Error(w, "AUTH_FAILURE", 500)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": dbUser.ID,
        "exp":     time.Now().Add(time.Hour * 12).Unix(),
    })
    
    tokenString, _ := token.SignedString(middleware.GetJWTKey())
    url := fmt.Sprintf("%s/auth-callback?token=%s&username=%s", os.Getenv("FRONTEND_URL"), tokenString, dbUser.Username)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleRegister initiates the registration process by staging the user 
// data and sending a 6-digit OTP email
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "INVALID_PAYLOAD", 400)
		return
	}

	// Check if a pending registration already exists
	var pending models.PendingUser
	err := db.DB.Where("email = ?", creds.Email).Find(&pending).Error

	if err == nil {
		if time.Now().Before(pending.ExpiresAt) {
			// resend existing code if it's still valid
			go sendEmail(pending.Email, pending.OTP)
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	otp := utils.GenerateOTP()

	newPending := models.PendingUser{
		Email:        creds.Email,
		Username:     creds.Username,
		PasswordHash: string(hashedPassword),
		OTP:          otp,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}

	// FirstOrCreate ensures we update any expired records for the same email
	if err := db.DB.Where(models.PendingUser{Email: creds.Email}).
		Assign(newPending).
		FirstOrCreate(&newPending).Error; err != nil {
		log.Printf("DB_ERROR: %v", err)
		http.Error(w, "SYSTEM_OVERLOAD", 500)
		return
	}

	// Send email in a goroutine
	go sendEmail(creds.Email, otp)
	w.WriteHeader(http.StatusAccepted)
}

// HandleVerify validates the 6-digit code. If correct, the user is 
// moved from 'PendingUser' to the main 'User' table.
func HandleVerify(w http.ResponseWriter, r *http.Request) {
	var req struct { Email string; Code string }
	json.NewDecoder(r.Body).Decode(&req)

	var pending models.PendingUser
	if err := db.DB.First(&pending, "email = ? AND otp = ?", req.Email, req.Code).Error; err != nil {
		http.Error(w, "Invalid code", 401)
		return
	}

	newUser := models.User{
		Username:     pending.Username,
		Email:        pending.Email,
		PasswordHash: pending.PasswordHash,
	}

	// Create real user and clean up pending record
	db.DB.Create(&newUser)
	db.DB.Delete(&pending)
	w.WriteHeader(http.StatusCreated)
}

// sendEmail connects to the configured SMTP server to deliver the OTP.
func sendEmail(to string, code string) {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	msg := fmt.Sprintf("Subject: GridGo Verification Code\r\n\r\nYour 6-digit access code is: %s\r\nExpires in 10 minutes.", code)

	auth := smtp.PlainAuth("", from, pass, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Email Failure: %v", err)
	}
}