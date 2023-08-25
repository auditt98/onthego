package validators

type IdPUserImportValidator struct {
	Id string `json:"id" binding:"required"`
}
