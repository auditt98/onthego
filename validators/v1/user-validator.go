package validators

type IdPUserImportValidator struct {
	Id       string `json:"id" binding:"required"`
	ClientId string `json:"clientId" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
}
