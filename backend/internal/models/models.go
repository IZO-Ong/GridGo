package models

import (
	"time"
)

type User struct {
    ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Username     string    `gorm:"uniqueIndex;not null"`
    Email        string    `gorm:"uniqueIndex;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time
    Mazes        []Maze    `gorm:"foreignKey:CreatorID"`
}

type PendingUser struct {
    Email        string    `gorm:"primaryKey"`
    Username     string    `gorm:"not null"`
    PasswordHash string    `gorm:"not null"`
    OTP          string    `gorm:"not null"`
    ExpiresAt    time.Time `gorm:"not null"`
}

type Maze struct {
    ID          string    `gorm:"primaryKey" json:"id"`
    CreatorID   *string   `gorm:"type:uuid" json:"creator_id"`
    WeightsJSON string    `gorm:"type:jsonb;not null" json:"weights_json"`
    Thumbnail   string    `gorm:"type:text" json:"thumbnail"`
    Rows        int       `gorm:"not null" json:"rows"`
    Cols        int       `gorm:"not null" json:"cols"`
    StartRow    int       `gorm:"not null" json:"start_row"`
    StartCol    int       `gorm:"not null" json:"start_col"`
    EndRow      int       `gorm:"not null" json:"end_row"`
    EndCol      int       `gorm:"not null" json:"end_col"`
    Complexity  float64   `json:"complexity"`
    DeadEnds    int       `json:"dead_ends"`
    CreatedAt   time.Time `json:"created_at"`
}