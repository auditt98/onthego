package validators

type NewAlbumValidator struct {
	Name string `json:"name" binding:"required"`
}

type AddUserToAlbumValidator struct {
	AlbumID uint     `json:"album_id" binding:"required"`
	UserIds []string `json:"user_ids" binding:"required"`
}
