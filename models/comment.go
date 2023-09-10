package models

import "time"

// User ...
type Comment struct {
	ID            uint      `gorm:"primarykey";json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Content       string    `json:"content"`
	AlbumID       uint      `json:"album_id"`
	Album         *Album    `json:"album"`
	PhotoID       uint      `json:"photo_id"`
	Photo         *Photo    `json:"photo"`
	CommenterID   string    `json:"commenter_id"`
	CommenterType string    `json:"commenter_type"`
}
