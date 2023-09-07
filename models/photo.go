package models

import "time"

type Photo struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	BaseName  string
	BaseUrl   string
	UserID    string
	User      *User
	AlbumID   uint
	Album     *Album
	// Sizes     []*ImageSize
	Likes    []*Like `gorm:"polymorphic:Liker;"`
	Comments []*Comment
}
