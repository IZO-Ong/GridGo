package models

import (
	"time"
)

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
    CreatedAt   time.Time `json:"created_at"`
}