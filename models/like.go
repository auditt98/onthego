package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	LikerID   uint
	LikerType string
	AlbumID   uint
	Album     *Album
	PhotoID   uint
	Photo     *Photo
}
