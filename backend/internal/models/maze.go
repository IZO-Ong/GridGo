// Package models defines the database schema and JSON structures for the application.
package models

import (
	"time"
)

// Maze represents a generated grid and its associated metadata.
type Maze struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	// CreatorID is a pointer to allow for maze generation by guest account.
	CreatorID   *string   `gorm:"type:uuid" json:"creator_id"`
	Creator     *User     `gorm:"foreignKey:CreatorID" json:"-"`
	// WeightsJSON stores the edge/wall weights of the grid in a optimized JSONB format.
	WeightsJSON string    `gorm:"type:jsonb;not null" json:"weights_json"`
	Thumbnail   string    `gorm:"type:text" json:"thumbnail"` // Base64 for gallery previews
	Rows        int       `gorm:"not null" json:"rows"`
	Cols        int       `gorm:"not null" json:"cols"`
	StartRow    int       `gorm:"not null" json:"start_row"`
	StartCol    int       `gorm:"not null" json:"start_col"`
	EndRow      int       `gorm:"not null" json:"end_row"`
	EndCol      int       `gorm:"not null" json:"end_col"`
	Complexity  float64   `json:"complexity"`
	CreatedAt   time.Time `json:"created_at"`
}