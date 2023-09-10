package models

import "time"

// User ...
type Album struct {
	ID         uint       `json:"id" gorm:"primaryKey" `
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Name       string     `json:"name"`
	Users      []*User    `gorm:"many2many:user_albums;" json:"users"`
	Photos     []*Photo   `json:"photos"`
	Likes      []*Like    `gorm:"polymorphic:Liker;" json:"likes"`
	Comments   []*Comment `gorm:"polymorphic:Commenter;" json:"comments"`
	LikesCount int        `gorm:"-" json:"likes_count"`
}
