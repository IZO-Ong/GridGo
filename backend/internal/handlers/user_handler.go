// Package handlers contains HTTP logic for the GridGo API.
// This file manages user profile data aggregation and social stats.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
)

// HandleGetProfile retrieves a comprehensive public profile for a specific username.
// It aggregates the user's mazes, forum posts and comments into a single response.
func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "USERNAME_REQUIRED", http.StatusBadRequest)
		return
	}

	// Check if person viewing the profile is logged in.
	viewerID := middleware.GetUserID(r)

	var user models.User
	// Perform a Preload to gather all associated content in a single DB hit.
	err := db.DB.Preload("Mazes").
		Preload("Posts.Creator").
		Preload("Comments.Creator").
		Preload("Comments.Post").
		Where("username = ?", username).
		First(&user).Error
	
	if err != nil {
		http.Error(w, "USER_NOT_FOUND", http.StatusNotFound)
		return
	}

	// If a viewer is logged in, we need to show them their own interaction status
	if viewerID != "" {
		// 1. Process Post Votes
		if len(user.Posts) > 0 {
			var postIDs []string
			for _, p := range user.Posts {
				postIDs = append(postIDs, p.ID)
			}

			var postVotes []models.PostVote
			db.DB.Where("user_id = ? AND post_id IN ?", viewerID, postIDs).Find(&postVotes)

			// Map vote values back to posts for O(1) lookup
			voteMap := make(map[string]int)
			for _, v := range postVotes {
				voteMap[v.PostID] = v.Value
			}

			for i := range user.Posts {
				user.Posts[i].UserVote = voteMap[user.Posts[i].ID]
			}
		}

		// Process Comment Votes
		if len(user.Comments) > 0 {
			var commentIDs []string
			for _, c := range user.Comments {
				commentIDs = append(commentIDs, c.ID)
			}

			var commentVotes []models.CommentVote
			db.DB.Where("user_id = ? AND comment_id IN ?", viewerID, commentIDs).Find(&commentVotes)

			voteMap := make(map[string]int)
			for _, v := range commentVotes {
				voteMap[v.CommentID] = v.Value
			}

			for i := range user.Comments {
				user.Comments[i].UserVote = voteMap[user.Comments[i].ID]
			}
		}
	}

	// Construct a JSON response object.
	response := map[string]interface{}{
		"username":   user.Username,
		"created_at": user.CreatedAt,
		"mazes":      user.Mazes,
		"posts":      user.Posts,
		"comments":   user.Comments,
		"stats": map[string]interface{}{
			"total_mazes":    len(user.Mazes),
			"total_posts":    len(user.Posts),
			"total_comments": len(user.Comments),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}