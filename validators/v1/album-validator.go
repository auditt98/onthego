package validators

type NewAlbumValidator struct {
	Name string `json:"name" binding:"required"`
}

type AddUserToAlbumValidator struct {
	UserIds []string `json:"user_ids" binding:"required"`
}
