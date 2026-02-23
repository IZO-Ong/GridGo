package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
)

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "USERNAME_REQUIRED", http.StatusBadRequest)
        return
    }

    viewerID := middleware.GetUserID(r)

    var user models.User
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

    if viewerID != "" {
        if len(user.Posts) > 0 {
            var postIDs []string
            for _, p := range user.Posts {
                postIDs = append(postIDs, p.ID)
            }

            var postVotes []models.PostVote
            db.DB.Where("user_id = ? AND post_id IN ?", viewerID, postIDs).Find(&postVotes)

            // Map vote values back to posts
            voteMap := make(map[string]int)
            for _, v := range postVotes {
                voteMap[v.PostID] = v.Value
            }

            for i := range user.Posts {
                user.Posts[i].UserVote = voteMap[user.Posts[i].ID]
            }
        }

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