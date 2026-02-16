package models

import (
	"time"
)

type User struct {
    ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Username     string    `gorm:"uniqueIndex;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time
    Mazes        []Maze    `gorm:"foreignKey:CreatorID"`
}

type Maze struct {
    ID          string    `gorm:"primaryKey"`
    CreatorID   *string   `gorm:"type:uuid"`
    GridJSON    string    `gorm:"type:jsonb;not null"`
    Rows        int       `gorm:"not null"`
    Cols        int       `gorm:"not null"`
    StartRow    int       `gorm:"not null"`
    StartCol    int       `gorm:"not null"`
    EndRow      int       `gorm:"not null"`
    EndCol      int       `gorm:"not null"`
    Complexity  float64
    DeadEnds    int
    CreatedAt   time.Time
}