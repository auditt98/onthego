package validators

type CommentValidator struct {
	Content string `json:"content" binding:"required"`
	AlbumID uint   `json:"album_id" binding:"required_without=PhotoID"`
	PhotoID uint   `json:"photo_id" binding:"required_without=AlbumID"`
}

type CommentUpdateValidator struct {
	Content string `json:"content" binding:"required"`
}
