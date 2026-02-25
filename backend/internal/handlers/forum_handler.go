// Package handlers contains the HTTP logic for the GridGo API.
// This file specifically manages forum interactions: Posts, Comments, and the Voting system.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/IZO-Ong/gridgo/internal/db"
	"github.com/IZO-Ong/gridgo/internal/middleware"
	"github.com/IZO-Ong/gridgo/internal/models"
	"github.com/IZO-Ong/gridgo/internal/utils"
)

// HandleCreatePost saves a new thread to the database.
// It generates a shortened custom ID (P-prefix) and links an optional MazeID.
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

	if p.MazeID != nil && *p.MazeID == "" {
		p.MazeID = nil
	}

	p.ID = utils.GeneratePostID()
	p.CreatorID = userID
	p.CreatedAt = time.Now()

	if err := db.DB.Create(&p).Error; err != nil {
		http.Error(w, "DB_ERROR", 500)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// HandleGetPosts retrieves a paginated list of posts for the main feed.
// It supports infinite scroll via the offset query parameter
func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	userID := middleware.GetUserID(r)
	
	var posts []models.Post
	// preload associations to avoid N+1 query problems
	db.DB.Preload("Creator").Preload("Maze").Preload("Comments").
		Order("created_at desc").Limit(10).Offset(offset).Find(&posts)

	// If logged in, fetch user specific votes for posts
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

		// inject user's vote value into Post struct for the frontend
		for i := range posts {
			posts[i].UserVote = voteMap[posts[i].ID]
		}
	}

	json.NewEncoder(w).Encode(posts)
}

// HandleGetPostByID fetches a single post with full detail, including its comments 
// and the creators of those comments.
func HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)
	
	var post models.Post

	// Preload
	err := db.DB.Preload("Creator").
				Preload("Maze"). 
				Preload("Comments.Creator").
				Where("id = ?", id).
				First(&post).Error
	
	if err != nil {
		http.Error(w, "POST_NOT_FOUND", 404)
		return
	}

	// Attach current user's vote status to the post and its comments
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

// HandleDeletePost removes a thread. It includes a check to ensure 
// the requester is the original creator.
func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete { return }
	
	postID := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)

	// Authorization check is built directly into the query
	result := db.DB.Where("id = ? AND creator_id = ?", postID, userID).Delete(&models.Post{})
	
	if result.RowsAffected == 0 {
		http.Error(w, "UNAUTHORIZED_OR_NOT_FOUND", 403)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleVote manages the logic for upvoting and downvoting.
// If a user votes the same way twice, the vote is removed (toggle behavior).
// If they change their vote, the record is updated.
func HandleVote(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "AUTH_REQUIRED", 401)
		return
	}

	var req struct {
		TargetID   string `json:"target_id"`   // ID of the post or comment
		TargetType string `json:"target_type"` // "post" or "comment"
		Value      int    `json:"value"`       // 1 for upvote, -1 for downvote
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "INVALID_INPUT", 400)
		return
	}

	if req.TargetType == "post" {
		var vote models.PostVote
		res := db.DB.Where("user_id = ? AND post_id = ?", userID, req.TargetID).First(&vote)
		if res.Error == nil {
			if vote.Value == req.Value { 
				db.DB.Delete(&vote) // Toggle off if same value
			} else {
				vote.Value = req.Value
				db.DB.Save(&vote)    // Update if different value
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

	// Trigger asynchronous recount of total votes for target
	updateVoteCount(req.TargetID, req.TargetType)
	w.WriteHeader(http.StatusOK)
}

// updateVoteCount is a helper function that recalculates the sum of votes
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

// HandleCreateComment adds a reply to a post.
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

	c.ID = utils.GenerateCommentID()
	c.CreatorID = userID
	c.CreatedAt = time.Now()

	if err := db.DB.Create(&c).Error; err != nil {
		http.Error(w, "DB_ERROR", 500)
		return
	}
	
	json.NewEncoder(w).Encode(c)
}

// HandleGetComments fetches all comments associated with a specific post, 
// sorted by vote count
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

// HandleDeleteComment removes a comment, ensuring only the owner can do so
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