package models

import "gorm.io/gorm"

// User ...
type Comment struct {
	gorm.Model
	Content       string
	AlbumID       uint
	Album         *Album
	PhotoID       uint
	Photo         *Photo
	CommenterID   uint
	CommenterType string
}
