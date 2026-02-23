package models

import (
	"time"
)

type Post struct {
    ID        string     `gorm:"primaryKey" json:"id"`
    Title     string     `json:"title"`
    Content   string     `json:"content"`
    MazeID    string     `json:"maze_id"` 
    Maze      Maze       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"maze"`
    CreatorID string     `json:"creator_id"`
    Creator   User       `gorm:"foreignKey:CreatorID" json:"creator"`
    Upvotes   int        `gorm:"default:0" json:"upvotes"`
    UserVote  int        `gorm:"-" json:"user_vote"`
    Comments  []Comment  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments"`
    Votes     []PostVote `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    CreatedAt time.Time  `json:"created_at"`
}

type Comment struct {
    ID        string        `gorm:"primaryKey" json:"id"`
    PostID    string        `json:"post_id"`
    Post      Post          `gorm:"foreignKey:PostID" json:"post"` 
    Content   string        `json:"content"`
    CreatorID string        `json:"creator_id"`
    Creator   User          `gorm:"foreignKey:CreatorID" json:"creator"`
    Upvotes   int           `gorm:"default:0" json:"upvotes"`
    UserVote  int           `gorm:"-" json:"user_vote"`
    Votes     []CommentVote `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
    CreatedAt time.Time     `json:"created_at"`
}

type PostVote struct {
    ID       uint   `gorm:"primaryKey"`
    UserID   string `gorm:"uniqueIndex:idx_post_user"`
    PostID   string `gorm:"uniqueIndex:idx_post_user"`
    Value    int    `json:"value"`
}

type CommentVote struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    string `gorm:"uniqueIndex:idx_comment_user"`
    CommentID string `gorm:"uniqueIndex:idx_comment_user"`
    Value     int    `json:"value"`
}