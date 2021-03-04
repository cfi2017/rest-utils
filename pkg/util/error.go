package util

type ErrorCode string

const (
	ErrorCodeBadRequest ErrorCode = "bad request"
	ErrorCodeNotFound   ErrorCode = "not found"
)

type ErrorResponse struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func NewErrorResponse(message string, code ErrorCode) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Code:    code,
	}
}
