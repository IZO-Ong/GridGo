package main

import (
	"log"
	"net/http"

	"github.com/IZO-Ong/gridgo/internal/auth"
	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/handlers"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func init() {
	godotenv.Load()
	key := middleware.GetJWTKey()
	store := sessions.NewCookieStore(key)
	store.Options = &sessions.Options{HttpOnly: true, SameSite: http.SameSiteLaxMode}
	gothic.Store = store
}

func main() {
	db.InitDB()
	auth.NewAuth()

	mux := http.NewServeMux()

    // Maze Routes
	mux.HandleFunc("/api/maze/generate", middleware.OptionalAuth(handlers.HandleGenerateMaze))
	mux.HandleFunc("/api/maze/get", handlers.HandleGetMaze)
	mux.HandleFunc("/api/maze/solve", handlers.HandleSolveMaze)
	mux.HandleFunc("/api/maze/render", handlers.HandleRenderMaze)
	mux.HandleFunc("/api/maze/thumbnail", handlers.HandleUpdateThumbnail)

    // User & Profile Routes
    mux.HandleFunc("/api/profile", handlers.HandleGetProfile)

    // Auth Routes
	mux.HandleFunc("/api/login", handlers.HandleLogin)
	mux.HandleFunc("/api/register", handlers.HandleRegister)
	mux.HandleFunc("/api/verify", handlers.HandleVerify)
	mux.HandleFunc("/api/auth/google", handlers.HandleOAuthLogin)
	mux.HandleFunc("/api/auth/google/callback", handlers.HandleOAuthCallback)

	log.Println("GridGo API online on :8080")
	http.ListenAndServe(":8080", middleware.EnableCORS(mux))
}