package models

import "time"

type Like struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LikerID   string    `json:"liker_id"`
	LikerType string    `json:"liker_type"`
	AlbumID   uint      `json:"album_id"`
	Album     *Album    `json:"album"`
	PhotoID   uint      `json:"photo_id"`
	Photo     *Photo    `json:"photo"`
}
