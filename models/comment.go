package models

import "time"

// User ...
type Comment struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Content       string
	AlbumID       uint
	Album         *Album
	PhotoID       uint
	Photo         *Photo
	CommenterID   uint
	CommenterType string
}
