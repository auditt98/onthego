package validators

type UpdatePhotoValidator struct {
	Name string `json:"name" binding:"required"`
}
