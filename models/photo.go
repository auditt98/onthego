package models

import "time"

type Photo struct {
	ID           uint      `gorm:"primarykey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	BaseName     string    `json:"base_name"`
	BaseUrl      string    `json:"base_url"`
	UserID       string    `json:"user_id"`
	User         *User     `json:"user"`
	AlbumID      uint      `json:"album_id"`
	Album        *Album    `json:"album"`
	PresignedUrl string    `json:"presigned_url"`
	// Sizes     []*ImageSize
	Likes    []*Like    `gorm:"polymorphic:Liker;" json:"likes"`
	Comments []*Comment `gorm:"polymorphic:Commenter;" json:"comments"`
}
