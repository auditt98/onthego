package models

import "time"

// User ...
type Album struct {
	ID         uint `json:"id" gorm:"primaryKey" `
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Users      []*User    `gorm:"many2many:user_albums;constraint:OnDelete:SET NULL;"`
	Photos     []*Photo   `gorm:"constraint:OnDelete:SET NULL;"`
	Likes      []*Like    `gorm:"polymorphic:Liker;constraint:OnDelete:SET NULL;"`
	Comments   []*Comment `gorm:"constraint:OnDelete:SET NULL;"`
	LikesCount int
}
