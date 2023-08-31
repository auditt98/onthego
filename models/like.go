package models

import "time"

type Like struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	LikerID   uint
	LikerType string
	AlbumID   uint
	Album     *Album
	PhotoID   uint
	Photo     *Photo
}
