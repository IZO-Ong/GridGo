package models

import (
	"time"
)

type User struct {
    ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Username     string    `gorm:"uniqueIndex;not null" json:"username"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    Mazes        []Maze    `gorm:"foreignKey:CreatorID" json:"mazes"`
    Posts        []Post    `gorm:"foreignKey:CreatorID" json:"posts"`
    Comments     []Comment `gorm:"foreignKey:CreatorID" json:"comments"`
    Votes        []Vote    `gorm:"foreignKey:UserID" json:"votes"`
}

type PendingUser struct {
    Email        string    `gorm:"primaryKey"`
    Username     string    `gorm:"not null"`
    PasswordHash string    `gorm:"not null"`
    OTP          string    `gorm:"not null"`
    ExpiresAt    time.Time `gorm:"not null"`
}