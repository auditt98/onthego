package validators

type AlbumValidator struct {
	Name string `json:"name" binding:"required"`
}
