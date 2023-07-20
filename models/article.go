package models

import (
	"gorm.io/gorm"
)

// Article ...
type Article struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
	User    User
}
