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

    // Maze Endpoints
	mux.HandleFunc("/api/maze/generate", middleware.OptionalAuth(handlers.HandleGenerateMaze))
	mux.HandleFunc("/api/maze/get", handlers.HandleGetMaze)
	mux.HandleFunc("/api/maze/my-mazes", middleware.RequireAuth(handlers.HandleGetMyMazes))
	mux.HandleFunc("/api/maze/delete", middleware.RequireAuth(handlers.HandleDeleteMaze))
	mux.HandleFunc("/api/maze/solve", handlers.HandleSolveMaze)
	mux.HandleFunc("/api/maze/render", handlers.HandleRenderMaze)
	mux.HandleFunc("/api/maze/thumbnail", handlers.HandleUpdateThumbnail)

    // User & Profile Endpoints
    mux.HandleFunc("/api/profile", handlers.HandleGetProfile)

	// Post Endpoints
	mux.HandleFunc("/api/forum/posts", middleware.OptionalAuth(handlers.HandleGetPosts))
	mux.HandleFunc("/api/forum/posts/create", middleware.RequireAuth(handlers.HandleCreatePost))
	mux.HandleFunc("/api/forum/post", handlers.HandleGetPostByID)
	mux.HandleFunc("/api/forum/post/delete", middleware.RequireAuth(handlers.HandleDeletePost))

	// Comment Endpoints
	mux.HandleFunc("/api/forum/comments", handlers.HandleGetComments)
	mux.HandleFunc("/api/forum/comment/create", middleware.RequireAuth(handlers.HandleCreateComment))
	mux.HandleFunc("/api/forum/comment/delete", middleware.RequireAuth(handlers.HandleDeleteComment))

	// Voting Endpoint
	mux.HandleFunc("/api/forum/vote", middleware.RequireAuth(handlers.HandleVote))

    // Auth Routes
	mux.HandleFunc("/api/login", handlers.HandleLogin)
	mux.HandleFunc("/api/register", handlers.HandleRegister)
	mux.HandleFunc("/api/verify", handlers.HandleVerify)
	mux.HandleFunc("/api/auth/google", handlers.HandleOAuthLogin)
	mux.HandleFunc("/api/auth/google/callback", handlers.HandleOAuthCallback)

	log.Println("GridGo API online on :8080")
	http.ListenAndServe(":8080", middleware.EnableCORS(mux))
}