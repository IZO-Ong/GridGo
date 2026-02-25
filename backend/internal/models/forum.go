// Package models defines the database schema and JSON structures for the application.
package models

import (
	"time"
)

// Post represents a forum thread, often used to share a specific Maze.
type Post struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	MazeID    *string    `json:"maze_id"` 
	// Maze association allows embedding the grid directly into the post view.
	Maze      Maze       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"maze"`
	CreatorID string     `json:"creator_id"`
	Creator   User       `gorm:"foreignKey:CreatorID" json:"creator"`
	// Upvotes is a denormalized field for fast sorting in the feed.
	Upvotes   int        `gorm:"default:0" json:"upvotes"`
	// UserVote is a virtual field (not in DB) used to tell the requester if they voted for this.
	UserVote  int        `gorm:"-" json:"user_vote"`
	Comments  []Comment  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments"`
	Votes     []PostVote `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}

// Comment represents a reply within a Post.
type Comment struct {
	ID        string        `gorm:"primaryKey" json:"id"`
	PostID    string        `json:"post_id"`
	Post      Post          `gorm:"foreignKey:PostID" json:"post"` 
	Content   string        `json:"content"`
	CreatorID string        `json:"creator_id"`
	Creator   User          `gorm:"foreignKey:CreatorID" json:"creator"`
	Upvotes   int           `gorm:"default:0" json:"upvotes"`
	UserVote  int           `gorm:"-" json:"user_vote"` // Virtual field for request context
	Votes     []CommentVote `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time     `json:"created_at"`
}

// PostVote handles the Many-to-Many relationship between Users and Posts.
// Value is 1 (up) or -1 (down).
type PostVote struct {
	ID      uint   `gorm:"primaryKey"`
	// uniqueIndex ensures a user can only have ONE vote per post.
	UserID  string `gorm:"uniqueIndex:idx_post_user"`
	PostID  string `gorm:"uniqueIndex:idx_post_user"`
	Value   int    `json:"value"`
}

// CommentVote handles the Many-to-Many relationship between Users and Comments.
type CommentVote struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    string `gorm:"uniqueIndex:idx_comment_user"`
	CommentID string `gorm:"uniqueIndex:idx_comment_user"`
	Value     int    `json:"value"`
}