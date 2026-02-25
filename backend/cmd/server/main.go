// Package main serves as the entry point for the GridGo API.
// It initializes the database, authentication providers, and defines the
// routing logic for maze management, forum interactions, and user authentication.
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

// init runs before main() to set up the environment and session management.
// It loads .env variables and configures the Gothic session store for OAuth2.
func init() {
	godotenv.Load()
	key := middleware.GetJWTKey()
	store := sessions.NewCookieStore(key)
	
	// Configure session options for security
	store.Options = &sessions.Options{
		HttpOnly: true, 
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = store
}

// main initializes the core services and starts the HTTP server.
func main() {
	// Initialize persistent services
	db.InitDB()    // Connects to the database
	auth.NewAuth() // Sets up OAuth providers

	// Create a new request multiplexer
	mux := http.NewServeMux()

	// --- Maze Endpoints ---
	// Handles the creation, retrieval, and visualization of maze grids.
	mux.HandleFunc("/api/maze/generate", middleware.OptionalAuth(handlers.HandleGenerateMaze))
	mux.HandleFunc("/api/maze/get", handlers.HandleGetMaze)
	mux.HandleFunc("/api/maze/my-mazes", middleware.RequireAuth(handlers.HandleGetMyMazes))
	mux.HandleFunc("/api/maze/delete", middleware.RequireAuth(handlers.HandleDeleteMaze))
	mux.HandleFunc("/api/maze/solve", handlers.HandleSolveMaze)
	mux.HandleFunc("/api/maze/render", handlers.HandleRenderMaze)
	mux.HandleFunc("/api/maze/thumbnail", handlers.HandleUpdateThumbnail)

	// --- User & Profile Endpoints ---
	// Manages user-specific metadata and settings.
	mux.HandleFunc("/api/profile", middleware.OptionalAuth(handlers.HandleGetProfile))

	// --- Post Endpoints ---
	// Core forum functionality for sharing mazes or general discussion.
	mux.HandleFunc("/api/forum/posts", middleware.OptionalAuth(handlers.HandleGetPosts))
	mux.HandleFunc("/api/forum/posts/create", middleware.RequireAuth(handlers.HandleCreatePost))
	mux.HandleFunc("/api/forum/post", middleware.OptionalAuth(handlers.HandleGetPostByID))
	mux.HandleFunc("/api/forum/post/delete", middleware.RequireAuth(handlers.HandleDeletePost))

	// --- Comment Endpoints ---
	// Handles nested interaction within forum posts.
	mux.HandleFunc("/api/forum/comments", middleware.OptionalAuth(handlers.HandleGetComments))
	mux.HandleFunc("/api/forum/comment/create", middleware.RequireAuth(handlers.HandleCreateComment))
	mux.HandleFunc("/api/forum/comment/delete", middleware.RequireAuth(handlers.HandleDeleteComment))

	// --- Voting Endpoint ---
	// Unified endpoint for upvoting/downvoting posts or comments.
	mux.HandleFunc("/api/forum/vote", middleware.RequireAuth(handlers.HandleVote))

	// --- Auth Routes ---
	// Manages standard email/password registration and Google OAuth flow.
	mux.HandleFunc("/api/login", handlers.HandleLogin)
	mux.HandleFunc("/api/register", handlers.HandleRegister)
	mux.HandleFunc("/api/verify", handlers.HandleVerify)
	mux.HandleFunc("/api/auth/google", handlers.HandleOAuthLogin)
	mux.HandleFunc("/api/auth/google/callback", handlers.HandleOAuthCallback)

	// Start server with CORS middleware
	log.Println("GridGo API online on :8080")
	http.ListenAndServe(":8080", middleware.EnableCORS(mux))
}