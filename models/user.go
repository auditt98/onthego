package models

// User ...
type User struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Email       string     `json:"email"`
	Albums      []*Album   `gorm:"many2many:user_albums; json:"albums"`
	LikedAlbums []*Album   `gorm:"many2many:album_likes;" json:"liked_albums"`
	LikedPhotos []*Photo   `gorm:"many2many:photo_likes;" json:"liked_photos"`
	Photos      []*Photo   `json:"photos"`
	Comments    []*Comment `gorm:"polymorphic:Commenter;" json:"comments"`
}
