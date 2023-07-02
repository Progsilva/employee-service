package handler

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewError(err error) *ErrorResponse {
	return &ErrorResponse{Message: err.Error()}
}
