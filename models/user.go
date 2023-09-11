package models

// User ...
type User struct {
	ID       string     `json:"id" gorm:"primaryKey"`
	Email    string     `json:"email"`
	Albums   []*Album   `gorm:"many2many:user_albums;" json:"albums"`
	Photos   []*Photo   `json:"photos"`
	Comments []*Comment `gorm:"polymorphic:Commenter;" json:"comments"`
	Likes    []*Like    `gorm:"polymorphic:Liker;" json:"likes"`
}
