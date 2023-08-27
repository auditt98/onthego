package models

// User ...
type User struct {
	Id string `json:"id" gorm:"primaryKey"`

	Albums      []*Album `gorm:"many2many:user_albums;"`
	LikedAlbums []*Album `gorm:"many2many:album_likes;"`
	LikedPhotos []*Photo `gorm:"many2many:photo_likes;"`
	Photos      []*Photo
	Comments    []*Comment `gorm:"polymorphic:Commenter;"`
}
