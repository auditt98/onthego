package validators

type NewAlbumValidator struct {
	Name string `json:"name" binding:"required"`
}

type AddUserToAlbumValidator struct {
	UserId string `json:"user_id" binding:"required"`
}
