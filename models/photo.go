package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	BaseName string
	UserID   uint
	User     *User
	AlbumID  uint
	Album    *Album
	Sizes    []*ImageSize
	Likes    []*Like `gorm:"polymorphic:Liker;"`
	Comments []*Comment
}
