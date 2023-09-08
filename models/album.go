package models

import "time"

// User ...
type Album struct {
	ID         uint `json:"id" gorm:"primaryKey" `
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Users      []*User `gorm:"many2many:user_albums;"`
	Photos     []*Photo
	Likes      []*Like `gorm:"polymorphic:Liker;"`
	Comments   []*Comment
	LikesCount int
}
