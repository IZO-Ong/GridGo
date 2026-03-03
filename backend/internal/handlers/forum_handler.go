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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- HELPERS & CACHE INVALIDATION ---

// invalidateFeedIndex clears the cached list of IDs for the feed.
// This must be called whenever a new post is created or one is deleted to ensure the order is correct.
func invalidateFeedIndex() {
	iter := db.RDB.Scan(db.Ctx, 0, "feed:ids:offset:*", 0).Iterator()
	for iter.Next(db.Ctx) {
		db.RDB.Del(db.Ctx, iter.Val())
	}
}

// attachVotesToPosts handles the slice version of your vote injector.
// It fetches user-specific vote data that should never be stored in the global Redis cache.
func attachVotesToPosts(userID string, posts []models.Post) {
	if userID == "" || len(posts) == 0 {
		return
	}

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

// attachVotesToPost fetches and injects the current user's vote status into a Post and its nested Comments.
func attachVotesToPost(userID string, post *models.Post) {
	if userID == "" || post == nil {
		return
	}

	// Fetch user's vote for main Post (using Find to avoid noise in logs)
	var postVote models.PostVote
	db.DB.Where("user_id = ? AND post_id = ?", userID, post.ID).Limit(1).Find(&postVote)
	post.UserVote = postVote.Value

	// Fetch user's votes for all nested Comments
	if len(post.Comments) > 0 {
		var commentIDs []string
		for _, c := range post.Comments {
			commentIDs = append(commentIDs, c.ID)
		}

		var cVotes []models.CommentVote
		db.DB.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&cVotes)

		voteMap := make(map[string]int)
		for _, v := range cVotes {
			voteMap[v.CommentID] = v.Value
		}

		for i := range post.Comments {
			post.Comments[i].UserVote = voteMap[post.Comments[i].ID]
		}
	}
}

// --- FORUM HANDLERS ---

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

	// Invalidate feed index because a new post shifts the order
	invalidateFeedIndex()

	json.NewEncoder(w).Encode(p)
}

// HandleGetPosts retrieves a paginated list of posts for the main feed.
// It supports infinite scroll via the offset query parameter and uses Two-Layer Caching.
func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetStr)
	userID := middleware.GetUserID(r)

	// Layer 1: Cache the list of IDs for this offset (The Index)
	idKey := "feed:ids:offset:" + offsetStr
	postIDs, err := db.GetOrSet(db.Ctx, idKey, 1*time.Hour, func() (*[]string, error) {
		var ids []string
		db.DB.Model(&models.Post{}).Order("created_at desc").Limit(10).Offset(offset).Pluck("id", &ids)
		return &ids, nil
	})

	if err != nil || postIDs == nil {
		http.Error(w, "DB_ERROR", 500)
		return
	}

	// Layer 2: Rehydrate the full Post objects from the Object Cache
	var finalPosts []models.Post
	for _, id := range *postIDs {
		cacheKey := "post:thread:" + id
		post, _ := db.GetOrSet(db.Ctx, cacheKey, 10*time.Minute, func() (*models.Post, error) {
			var p models.Post
			err := db.DB.Preload("Creator").Preload("Maze").Preload("Comments.Creator").
				Where("id = ?", id).First(&p).Error
			return &p, err
		})
		if post != nil {
			finalPosts = append(finalPosts, *post)
		}
	}

	// Inject personal vote data (Never cached globally)
	if userID != "" && len(finalPosts) > 0 {
		attachVotesToPosts(userID, finalPosts)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalPosts)
}

// HandleGetPostByID fetches a single post with full detail, including its comments 
// and the creators of those comments.
func HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)
	cacheKey := "post:thread:" + id

	post, err := db.GetOrSet(db.Ctx, cacheKey, 10*time.Minute, func() (*models.Post, error) {
		var p models.Post
		err := db.DB.Preload("Creator").Preload("Maze").Preload("Comments.Creator").
			Where("id = ?", id).First(&p).Error
		return &p, err
	})

	if err != nil {
		http.Error(w, "POST_NOT_FOUND", 404)
		return
	}

	if userID != "" {
		attachVotesToPost(userID, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// HandleDeletePost removes a thread. It includes a check to ensure 
// the requester is the original creator.
func HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		return
	}

	postID := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)

	result := db.DB.Where("id = ? AND creator_id = ?", postID, userID).Delete(&models.Post{})

	if result.RowsAffected == 0 {
		http.Error(w, "UNAUTHORIZED_OR_NOT_FOUND", 403)
		return
	}

	// Clear the object from cache
	db.RDB.Del(db.Ctx, "post:thread:"+postID)
	// Clear the feed index because the order has shifted
	invalidateFeedIndex()

	w.WriteHeader(http.StatusOK)
}

// HandleVote manages the logic for upvoting and downvoting.
// This version removes atomic increments in favor of direct database updates 
// and cache invalidation to ensure 100% data consistency.
func HandleVote(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "AUTH_REQUIRED", 401)
		return
	}

	var req struct {
		TargetID   string `json:"target_id"`   // ID of the post or comment
		TargetType string `json:"target_type"` // "post" or "comment"
		Value      int    `json:"value"`       // Target: 1 (up), -1 (down), or 0 (remove)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "INVALID_INPUT", 400)
		return
	}

	// We use a transaction to ensure the vote update and the count recount 
	// happen together atomically.
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if req.TargetType == "post" {
			if req.Value == 0 {
				// Remove vote record if it exists
				tx.Where("user_id = ? AND post_id = ?", userID, req.TargetID).Delete(&models.PostVote{})
			} else {
				// Upsert vote record to handle race conditions without unique constraint errors
				err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"value"}),
				}).Create(&models.PostVote{
					UserID: userID,
					PostID: req.TargetID,
					Value:  req.Value,
				}).Error
				if err != nil {
					return err
				}
			}

			// Calculate new sum and update Posts table directly
			var total int64
			tx.Model(&models.PostVote{}).Where("post_id = ?", req.TargetID).Select("COALESCE(SUM(value), 0)").Row().Scan(&total)
			tx.Model(&models.Post{}).Where("id = ?", req.TargetID).Update("upvotes", total)

			// Clear the post thread cache
			db.RDB.Del(db.Ctx, "post:thread:"+req.TargetID)

		} else {
			if req.Value == 0 {
				tx.Where("user_id = ? AND comment_id = ?", userID, req.TargetID).Delete(&models.CommentVote{})
			} else {
				err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_id"}, {Name: "comment_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"value"}),
				}).Create(&models.CommentVote{
					UserID:    userID,
					CommentID: req.TargetID,
					Value:     req.Value,
				}).Error
				if err != nil {
					return err
				}
			}

			var total int64
			tx.Model(&models.CommentVote{}).Where("comment_id = ?", req.TargetID).Select("COALESCE(SUM(value), 0)").Row().Scan(&total)
			tx.Model(&models.Comment{}).Where("id = ?", req.TargetID).Update("upvotes", total)

			var c models.Comment
			if err := tx.Select("post_id").First(&c, "id = ?", req.TargetID).Error; err == nil {
				db.RDB.Del(db.Ctx, "post:thread:"+c.PostID)
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, "DB_ERROR", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	// Clear parent post object so new comment is shown in detail view
	db.RDB.Del(db.Ctx, "post:thread:"+c.PostID)

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
	if r.Method != http.MethodDelete {
		return
	}

	commentID := r.URL.Query().Get("id")
	userID := middleware.GetUserID(r)

	// Fetch parent post ID before deleting to clear cache
	var c models.Comment
	db.DB.Select("post_id").Where("id = ?", commentID).First(&c)

	result := db.DB.Where("id = ? AND creator_id = ?", commentID, userID).Delete(&models.Comment{})

	if result.RowsAffected == 0 {
		http.Error(w, "UNAUTHORIZED", 403)
		return
	}

	// Invalidate the parent post thread cache
	db.RDB.Del(db.Ctx, "post:thread:"+c.PostID)

	w.WriteHeader(http.StatusOK)
}