// Package models defines the database schema and JSON structures for the application.
package models

import (
	"time"
)

// User represents a registered account.
// It serves as the primary entity for mazes, forum participation and social interactions.
type User struct {
	// ID uses UUID for global uniqueness and security against sequential ID scraping.
	ID           string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string        `gorm:"uniqueIndex;not null" json:"username"`
	Email        string        `gorm:"uniqueIndex;not null" json:"email"`
	// PasswordHash is marked with json:"-" to ensure it is not leaked in API responses.
	PasswordHash string        `gorm:"not null" json:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	
	// Relationships: One-to-Many
	Mazes        []Maze        `gorm:"foreignKey:CreatorID" json:"mazes"`
	Posts        []Post        `gorm:"foreignKey:CreatorID" json:"posts"`
	Comments     []Comment     `gorm:"foreignKey:CreatorID" json:"comments"`
	PostVotes    []PostVote    `gorm:"foreignKey:UserID" json:"post_votes"`
	CommentVotes []CommentVote `gorm:"foreignKey:UserID" json:"comment_votes"`
}

// PendingUser acts as a staging area for registration.
// Records are held here until the OTP is verified.
type PendingUser struct {
	Email        string    `gorm:"primaryKey"` // Email is the unique identifier for verification
	Username     string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null"`
	OTP          string    `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`
}