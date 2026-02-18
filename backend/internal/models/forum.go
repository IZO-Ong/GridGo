package models

import (
	"time"
)

type Post struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    MazeID    string    `json:"maze_id"`
    Maze      Maze      `gorm:"foreignKey:MazeID" json:"maze"`
    CreatorID string    `json:"creator_id"`
    Creator   User      `gorm:"foreignKey:CreatorID" json:"creator"`
    Upvotes   int       `gorm:"default:0" json:"upvotes"`
    UserVote  int       `gorm:"-" json:"user_vote"`
    Comments  []Comment `json:"comments"`
    CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	PostID    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatorID string    `json:"creator_id"`
	Creator   User      `gorm:"foreignKey:CreatorID" json:"creator"`
	Upvotes   int       `gorm:"default:0" json:"upvotes"`
	CreatedAt time.Time `json:"created_at"`
}

type Vote struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     string `gorm:"uniqueIndex:idx_user_target"`
	TargetID   string `gorm:"uniqueIndex:idx_user_target"`
	TargetType string `json:"target_type"`
	Value      int    `json:"value"`
}