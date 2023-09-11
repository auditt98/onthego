package validators

type LikeValidator struct {
	AlbumID uint `json:"album_id" binding:"required_without=PhotoID"`
	PhotoID uint `json:"photo_id" binding:"required_without=AlbumID"`
}
