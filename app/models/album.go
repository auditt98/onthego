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
	Likes      []*Like    `json:"likes"`
	Comments   []*Comment `json:"comments"`
	LikesCount int        `json:"likes_count"`
}
