package models

import "time"

type ImageSize struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	SizeName  string
	Width     uint
	Height    uint
	URL       string
	PhotoID   uint
	Photo     *Photo
}
