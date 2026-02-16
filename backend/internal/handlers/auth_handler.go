package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Identifier string `json:"username"` 
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
    if err := db.DB.Where("username = ? OR email = ?", creds.Identifier, creds.Identifier).First(&user).Error; err != nil {
        http.Error(w, "INVALID_CREDENTIALS", http.StatusUnauthorized)
        return
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(middleware.GetJWTKey())
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "username": user.Username})
}

func HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    q.Set("provider", "google")
    r.URL.RawQuery = q.Encode()
    
    gothic.BeginAuthHandler(w, r)
}

func HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil { return }

	var dbUser models.User
	if db.DB.Where("email = ?", user.Email).First(&dbUser).Error != nil {
		dbUser = models.User{Username: user.NickName, Email: user.Email, PasswordHash: "OAUTH_ACCOUNT"}
		db.DB.Create(&dbUser)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(middleware.GetJWTKey())
	url := fmt.Sprintf("%s/auth-callback?token=%s&username=%s", os.Getenv("FRONTEND_URL"), tokenString, dbUser.Username)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"email"`
		Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "INVALID_PAYLOAD", http.StatusBadRequest)
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
    otp := GenerateOTP()

    pending := models.PendingUser{
        Email:        creds.Email,
        Username:     creds.Username,
        PasswordHash: string(hashedPassword),
        OTP:          otp,
        ExpiresAt:    time.Now().Add(10 * time.Minute),
    }

    if err := db.DB.Save(&pending).Error; err != nil {
        http.Error(w, "Email already in use for registration", http.StatusConflict)
        return
    }

    go sendEmail(creds.Email, otp)
    w.WriteHeader(http.StatusAccepted)
}

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

    db.DB.Create(&newUser)
    db.DB.Delete(&pending)
    w.WriteHeader(http.StatusCreated)
}

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