package models

import "time"

type ImageSize struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SizeName  string    `json:"size_name"`
	Width     uint      `json:"width"`
	Height    uint      `json:"height"`
	URL       string    `json:"url"`
	PhotoID   uint      `json:"photo_id"`
	Photo     *Photo    `json:"photo"`
}
