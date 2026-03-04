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

// HELPERS & CACHE INVALIDATION

// invalidateFeedIndex clears the cached list of IDs for the feed.
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

	// Fetch user's vote for main Post
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

// FORUM HANDLERS

// HandleCreatePost saves a new thread to the database.
// It generates a shortened custom ID and links an optional MazeID.
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

    invalidateFeedIndex()
    
    var username string
    if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Pluck("username", &username).Error; err == nil {
        db.RDB.Del(db.Ctx, "profile:public:"+username)
    }

    json.NewEncoder(w).Encode(p)
}

// HandleGetPosts retrieves a paginated list of posts for the main feed.
func HandleGetPosts(w http.ResponseWriter, r *http.Request) {
    offsetStr := r.URL.Query().Get("offset")
    offset, _ := strconv.Atoi(offsetStr)
    userID := middleware.GetUserID(r)

    // Fetch IDs
    idKey := "feed:ids:offset:" + offsetStr
    postIDs, err := db.GetOrSet(db.Ctx, idKey, 1*time.Hour, func() (*[]string, error) {
        var ids []string
        err := db.DB.Model(&models.Post{}).Order("created_at desc").Limit(10).Offset(offset).Pluck("id", &ids).Error
        return &ids, err
    })

    if err != nil || postIDs == nil {
        http.Error(w, "DB_ERROR", 500)
        return
    }

    // Initialize as an empty slice
    finalPosts := make([]models.Post, 0) 

    // Rehydrate objects
    for _, id := range *postIDs {
        cacheKey := "post:thread:" + id
        post, _ := db.GetOrSet(db.Ctx, cacheKey, 10*time.Minute, func() (*models.Post, error) {
            var p models.Post
            // FIX 2: Check MazeID before preloading to avoid rehydration errors
            query := db.DB.Preload("Creator").Preload("Comments.Creator")
            
            // Only Preload Maze if the post actually has one linked
            err := query.Where("id = ?", id).First(&p).Error
            if err == nil && p.MazeID != nil {
                db.DB.Model(&p).Association("Maze").Find(&p.Maze)
            }
            
            return &p, err
        })
        
        if post != nil {
            finalPosts = append(finalPosts, *post)
        }
    }

    // Inject personal votes
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
    if r.Method != http.MethodDelete { return }

    postID := r.URL.Query().Get("id")
    userID := middleware.GetUserID(r)

    var post models.Post
    if err := db.DB.Preload("Creator").Where("id = ? AND creator_id = ?", postID, userID).First(&post).Error; err != nil {
        http.Error(w, "NOT_FOUND_OR_UNAUTHORIZED", 403)
        return
    }

    // Perform the delete
    if err := db.DB.Delete(&post).Error; err != nil {
        http.Error(w, "DB_ERROR", 500)
        return
    }

    // Clear thread, feed index, and the USER PROFILE
    db.RDB.Del(db.Ctx, "post:thread:"+postID)
    db.RDB.Del(db.Ctx, "profile:public:"+post.Creator.Username)
    invalidateFeedIndex()

    w.WriteHeader(http.StatusOK)
}

// HandleVote manages the logic for upvoting and downvoting.
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

    err := db.DB.Transaction(func(tx *gorm.DB) error {
        var targetUsername string

        if req.TargetType == "post" {
            if req.Value == 0 {
                tx.Where("user_id = ? AND post_id = ?", userID, req.TargetID).Delete(&models.PostVote{})
            } else {
                tx.Clauses(clause.OnConflict{
                    Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
                    DoUpdates: clause.AssignmentColumns([]string{"value"}),
                }).Create(&models.PostVote{UserID: userID, PostID: req.TargetID, Value: req.Value})
            }

            var total int64
            tx.Model(&models.PostVote{}).Where("post_id = ?", req.TargetID).Select("COALESCE(SUM(value), 0)").Row().Scan(&total)
            tx.Model(&models.Post{}).Where("id = ?", req.TargetID).Update("upvotes", total)

            tx.Table("posts").
                Joins("JOIN users ON users.id = posts.creator_id").
                Where("posts.id = ?", req.TargetID).
                Pluck("users.username", &targetUsername)

            db.RDB.Del(db.Ctx, "post:thread:"+req.TargetID)

        } else {
            if req.Value == 0 {
                tx.Where("user_id = ? AND comment_id = ?", userID, req.TargetID).Delete(&models.CommentVote{})
            } else {
                tx.Clauses(clause.OnConflict{
                    Columns:   []clause.Column{{Name: "user_id"}, {Name: "comment_id"}},
                    DoUpdates: clause.AssignmentColumns([]string{"value"}),
                }).Create(&models.CommentVote{UserID: userID, CommentID: req.TargetID, Value: req.Value})
            }

            var total int64
            tx.Model(&models.CommentVote{}).Where("comment_id = ?", req.TargetID).Select("COALESCE(SUM(value), 0)").Row().Scan(&total)
            tx.Model(&models.Comment{}).Where("id = ?", req.TargetID).Update("upvotes", total)

            var c models.Comment
            if err := tx.Preload("Creator").First(&c, "id = ?", req.TargetID).Error; err == nil {
                targetUsername = c.Creator.Username
                db.RDB.Del(db.Ctx, "post:thread:"+c.PostID)
            }
        }

        if targetUsername != "" {
            db.RDB.Del(db.Ctx, "profile:public:"+targetUsername)
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

    db.RDB.Del(db.Ctx, "post:thread:"+c.PostID)

    var username string
    if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Pluck("username", &username).Error; err == nil {
        db.RDB.Del(db.Ctx, "profile:public:"+username)
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

    var c models.Comment
    if err := db.DB.Preload("Creator").Where("id = ? AND creator_id = ?", commentID, userID).First(&c).Error; err != nil {
        http.Error(w, "NOT_FOUND_OR_UNAUTHORIZED", 403)
        return
    }

    db.DB.Delete(&c)

    db.RDB.Del(db.Ctx, "post:thread:"+c.PostID)
    db.RDB.Del(db.Ctx, "profile:public:"+c.Creator.Username)

    w.WriteHeader(http.StatusOK)
}