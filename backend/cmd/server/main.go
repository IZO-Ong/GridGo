package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/IZO-Ong/gridgo/internal/auth"
	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/maze"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set in environment")
	}
	return []byte(secret)
}

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found")
    }

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is required for session store")
    }
    
    key := []byte(secret) 
    store := sessions.NewCookieStore(key)
    
    store.Options.HttpOnly = true
    store.Options.Secure = false
    store.Options.SameSite = http.SameSiteLaxMode
    
    gothic.Store = store
}

func main() {

	db.InitDB()
	auth.NewAuth()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/maze/get", handleGetMaze)
	mux.HandleFunc("/api/maze/generate", handleGenerateMaze)
	mux.HandleFunc("/api/maze/render", handleRenderMaze)
	mux.HandleFunc("/api/maze/solve", handleSolveMaze)
    mux.HandleFunc("/api/register", handleRegister)
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/verify", handleVerify)
	mux.HandleFunc("/api/auth/google", handleOAuthLogin)
    mux.HandleFunc("/api/auth/google/callback", handleOAuthCallback)

	println("GridGo API running on port 8080")
	http.ListenAndServe(":8080", enableCORS(mux))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL")) 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleGenerateMaze(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // parse metadata
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    rows, _ := strconv.Atoi(r.FormValue("rows"))
    cols, _ := strconv.Atoi(r.FormValue("cols"))

    // check boundary
    if rows < 2 || rows > 300 || cols < 2 || cols > 300 {
        http.Error(w, "OUT_OF_BOUNDS: Dimensions must be between 2 and 300", http.StatusBadRequest)
        return
    }

    genType := r.FormValue("type")
    myMaze := maze.NewMaze(rows, cols)

	switch genType {
	case "image":
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image required for image-type maze", http.StatusBadRequest)
			return
		}
		defer file.Close()

		weights, err := maze.GetEdgeWeights(file, rows, cols)
		if err != nil {
			http.Error(w, "Vision processing failed", http.StatusInternalServerError)
			return
		}
		myMaze.GenerateImageMaze(weights)

	case "kruskal":
		myMaze.GenerateKruskal()

	case "recursive":
		myMaze.GenerateRecursive(0, 0)

	default:
		http.Error(w, "Invalid generation type", http.StatusBadRequest)
		return
	}

    myMaze.SetRandomStartEnd()
	weightsBytes, _ := json.Marshal(myMaze.Weights)
    stats := myMaze.CalculateStats()

    mazeID := "M-" + strconv.Itoa(rand.Intn(9000)+1000) + "-X"

    dbMaze := models.Maze{
        ID:           mazeID,
		WeightsJSON:  string(weightsBytes),
        Rows:         rows,
        Cols:         cols,
        StartRow:     myMaze.Start[0],
        StartCol:     myMaze.Start[1],
        EndRow:       myMaze.End[0],
        EndCol:       myMaze.End[1],
        DeadEnds:     stats.DeadEnds,
        Complexity:   stats.Complexity,
    }

    myMaze.ID = mazeID
    myMaze.DeadEnds = stats.DeadEnds
    myMaze.Complexity = stats.Complexity

    result := db.DB.Create(&dbMaze)
    if result.Error != nil {
        http.Error(w, "Failed to save to database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(myMaze)
}

func handleRenderMaze(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions { return }

    var m maze.Maze
    if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
        http.Error(w, "Invalid data", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "image/png")
    m.RenderToWriter(w, 10) 
}

func handleSolveMaze(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions { return }
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Maze      maze.Maze `json:"maze"`
		Algorithm string    `json:"algorithm"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var visited [][2]int
	var path [][2]int

	switch payload.Algorithm {
	case "astar":
		visited, path = payload.Maze.SolveAStar()
	case "bfs":
		visited, path = payload.Maze.SolveBFS()
	case "greedy":
    	visited, path = payload.Maze.SolveGreedy()
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"visited": visited,
		"path":    path,
	})
}

func handleGetMaze(w http.ResponseWriter, r *http.Request) {
    mazeID := r.URL.Query().Get("id")
    if mazeID == "" {
        http.Error(w, "Maze ID required", http.StatusBadRequest)
        return
    }

    var m models.Maze
    if err := db.DB.First(&m, "id = ?", mazeID).Error; err != nil {
        http.Error(w, "Maze not found", http.StatusNotFound)
        return
    }

    reconstructed := maze.NewMaze(m.Rows, m.Cols)

    var savedWeights map[string]int
    json.Unmarshal([]byte(m.WeightsJSON), &savedWeights)

    reconstructed.GenerateImageMaze(savedWeights)

    reconstructed.SetManualStartEnd(m.StartRow, m.StartCol, m.EndRow, m.EndCol)

    response := map[string]interface{}{
        "id":         m.ID,
        "rows":       m.Rows,
        "cols":       m.Cols,
        "grid":       reconstructed.Grid, // Re-generated walls
        "start":      [2]int{m.StartRow, m.StartCol},
        "end":        [2]int{m.EndRow, m.EndCol},
        "dead_ends":  m.DeadEnds,
        "complexity": m.Complexity,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	var user models.User
	if err := db.DB.First(&user, "username = ?", creds.Username).Error; err != nil {
		http.Error(w, "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, "INVALID_CREDENTIALS", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString(getJWTKey())
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "username": user.Username})
}

func GenerateOTP() string {
    // Standard 6-digit numeric code
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email    string `json:"username"`
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
        PasswordHash: string(hashedPassword),
        OTP:          otp,
        ExpiresAt:    time.Now().Add(10 * time.Minute),
    }

    if err := db.DB.Save(&pending).Error; err != nil {
        http.Error(w, "DATABASE_ERROR", http.StatusInternalServerError)
        return
    }

    go sendEmail(creds.Email, otp)
    
    w.WriteHeader(http.StatusAccepted)
}

// Final account creation after 6-digit code check
func handleVerify(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Code  string `json:"code"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    var pending models.PendingUser
    if err := db.DB.First(&pending, "email = ? AND otp = ?", req.Email, req.Code).Error; err != nil {
        http.Error(w, "INVALID_OR_EXPIRED_CODE", http.StatusUnauthorized)
        return
    }

    if time.Now().After(pending.ExpiresAt) {
        db.DB.Delete(&pending)
        http.Error(w, "CODE_EXPIRED", http.StatusUnauthorized)
        return
    }

    user := models.User{
        Username:     pending.Email, 
        PasswordHash: pending.PasswordHash,
    }

    if err := db.DB.Create(&user).Error; err != nil {
        http.Error(w, "IDENTITY_COLLISION", http.StatusConflict)
        return
    }

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

func handleOAuthLogin(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    q.Set("provider", "google")
    r.URL.RawQuery = q.Encode()
    
    gothic.BeginAuthHandler(w, r)
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
    q.Set("provider", "google")
    r.URL.RawQuery = q.Encode()

    user, err := gothic.CompleteUserAuth(w, r)
    frontendURL := os.Getenv("FRONTEND_URL")

    if err != nil {
        http.Redirect(w, r, frontendURL+"/login?error=OAUTH_FAILED", http.StatusTemporaryRedirect)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.Email,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })
    
    tokenString, _ := token.SignedString(getJWTKey())

    url := fmt.Sprintf("%s/auth-callback?token=%s&username=%s", frontendURL, tokenString, user.NickName)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}