package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/google/uuid"
)

// HandleCreatePost handles new forum threads
func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "UNAUTHORIZED", 401)
		return
	}

	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "INVALID_PAYLOAD", 400)
		return
	}

	p.ID = "P-" + uuid.New().String()[:8]
	p.CreatorID = userID
	p.CreatedAt = time.Now()

	if err := db.DB.Create(&p).Error; err != nil {
		http.Error(w, "DB_ERROR", 500)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// HandleGetPosts supports infinite scroll via offset
func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    userID := middleware.GetUserID(r)
    
    var posts []models.Post
    db.DB.Preload("Creator").Preload("Maze").Preload("Comments").Order("created_at desc").Limit(10).Offset(offset).Find(&posts)

    if userID != "" && len(posts) > 0 {
        var postIDs []string
        for _, p := range posts {
            postIDs = append(postIDs, p.ID) 
        }

        var votes []models.PostVote
        db.DB.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&votes)

        voteMap := make(map[string]int)
        for _, v := range votes {
            voteMap[v.PostID] = v.Value
        }

        for i := range posts {
            posts[i].UserVote = voteMap[posts[i].ID]
        }
    }

    json.NewEncoder(w).Encode(posts)
}

// HandleGetPostByID fetches a post and its comments
func HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    userID := middleware.GetUserID(r)

	fmt.Printf("API_CALL: GetPostByID | ID: '%s' | User: '%s'\n", id, userID)
	
    var post models.Post

    err := db.DB.Preload("Creator").
                Preload("Maze"). 
                Preload("Comments.Creator").
                Where("id = ?", id).
                First(&post).Error
    
    if err != nil {
        http.Error(w, "POST_NOT_FOUND", 404)
        return
    }

    if userID != "" {
		var postVote models.PostVote
		db.DB.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&postVote)
		post.UserVote = postVote.Value

		if len(post.Comments) > 0 {
			var commentIDs []string
			for _, c := range post.Comments { commentIDs = append(commentIDs, c.ID) }

			var cVotes []models.CommentVote
			db.DB.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&cVotes)

			voteMap := make(map[string]int)
			for _, v := range cVotes { voteMap[v.CommentID] = v.Value }
			for i := range post.Comments {
				post.Comments[i].UserVote = voteMap[post.Comments[i].ID]
			}
		}
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// HandleDeletePost ensures only the author can delete the thread
func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete { return }
	
	postID := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)

	result := db.DB.Where("id = ? AND creator_id = ?", postID, userID).Delete(&models.Post{})
	
	if result.RowsAffected == 0 {
		http.Error(w, "UNAUTHORIZED_OR_NOT_FOUND", 403)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleVote manages the upvote/downvote toggle
func HandleVote(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    if userID == "" {
        http.Error(w, "AUTH_REQUIRED", 401)
        return
    }

    var req struct {
        TargetID   string `json:"target_id"`
        TargetType string `json:"target_type"`
        Value      int    `json:"value"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "INVALID_INPUT", 400)
        return
    }

    if req.TargetType == "post" {
        var vote models.PostVote
        res := db.DB.Where("user_id = ? AND post_id = ?", userID, req.TargetID).First(&vote)
        if res.Error == nil {
            if vote.Value == req.Value { db.DB.Delete(&vote) } else {
                vote.Value = req.Value
                db.DB.Save(&vote)
            }
        } else {
            db.DB.Create(&models.PostVote{UserID: userID, PostID: req.TargetID, Value: req.Value})
        }
    } else {
        var vote models.CommentVote
        res := db.DB.Where("user_id = ? AND comment_id = ?", userID, req.TargetID).First(&vote)
        if res.Error == nil {
            if vote.Value == req.Value { db.DB.Delete(&vote) } else {
                vote.Value = req.Value
                db.DB.Save(&vote)
            }
        } else {
            db.DB.Create(&models.CommentVote{UserID: userID, CommentID: req.TargetID, Value: req.Value})
        }
    }

    updateVoteCount(req.TargetID, req.TargetType)
    w.WriteHeader(http.StatusOK)
}

func updateVoteCount(targetID, targetType string) {
    var total int64
    if targetType == "post" {
        db.DB.Model(&models.PostVote{}).Where("post_id = ?", targetID).Select("SUM(value)").Row().Scan(&total)
        db.DB.Model(&models.Post{}).Where("id = ?", targetID).Update("upvotes", total)
    } else {
        db.DB.Model(&models.CommentVote{}).Where("comment_id = ?", targetID).Select("SUM(value)").Row().Scan(&total)
        db.DB.Model(&models.Comment{}).Where("id = ?", targetID).Update("upvotes", total)
    }
}

// HandleCreateComment adds a new flat comment to a post
func HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "AUTH_REQUIRED", 401)
		return
	}

	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "INVALID_PAYLOAD", 400)
		return
	}

	c.ID = fmt.Sprintf("C-%d", rand.Intn(1000000))
	c.CreatorID = userID
	c.CreatedAt = time.Now()

	if err := db.DB.Create(&c).Error; err != nil {
		http.Error(w, "DB_ERROR", 500)
		return
	}
	json.NewEncoder(w).Encode(c)
}

// HandleGetComments fetches comments for a post, sorted by upvotes
func HandleGetComments(w http.ResponseWriter, r *http.Request) {
    postID := r.URL.Query().Get("post_id")
    userID := middleware.GetUserID(r)
    var comments []models.Comment

    db.DB.Preload("Creator").Where("post_id = ?", postID).Order("upvotes desc").Find(&comments)

    if userID != "" && len(comments) > 0 {
        var commentIDs []string
        for _, c := range comments {
            commentIDs = append(commentIDs, c.ID)
        }

        // Query the specific CommentVote table
        var votes []models.CommentVote
        db.DB.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&votes)

        voteMap := make(map[string]int)
        for _, v := range votes {
            voteMap[v.CommentID] = v.Value
        }

        for i := range comments {
            comments[i].UserVote = voteMap[comments[i].ID]
        }
    }

    json.NewEncoder(w).Encode(comments)
}

// HandleDeleteComment restricts deletion to the comment owner
func HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete { return }
	
	commentID := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)

	result := db.DB.Where("id = ? AND creator_id = ?", commentID, userID).Delete(&models.Comment{})
	
	if result.RowsAffected == 0 {
		http.Error(w, "UNAUTHORIZED", 403)
		return
	}
	w.WriteHeader(http.StatusOK)
}