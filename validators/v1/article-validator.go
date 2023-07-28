package validators

type ArticleValidator struct {
	Name string `json:"name" form:"article_name" binding:"required"`
	User int32  `json:"user" form:"article_user" binding:"required"`
}
