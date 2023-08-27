package models

import "time"

// User ...
type Album struct {
	ID uint `json:"id" gorm:"primaryKey" `

	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []*User `gorm:"many2many:album_users;"`

	Likes      []*User `gorm:"many2many:album_likes;"`
	LikesCount int
	Comments   []*Comment
	Photos     []Photo
}
