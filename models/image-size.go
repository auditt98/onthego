package models

import "gorm.io/gorm"

type ImageSize struct {
	gorm.Model
	SizeName string
	Width    uint
	Height   uint
	URL      string
	PhotoID  uint
	Photo    *Photo
}
