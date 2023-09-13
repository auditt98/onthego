package types

type SuccessSearchResponse struct {
	Data     interface{} `json:"data"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
