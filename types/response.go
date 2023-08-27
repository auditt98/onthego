package types

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
