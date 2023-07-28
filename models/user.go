package models

import (
	"gorm.io/gorm"
)

// User ...
type User struct {
	gorm.Model
	Email    string
	Password string
	Name     string
	Articles []Article
}
