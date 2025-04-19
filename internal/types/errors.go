package types

type ErrorCode string

const (
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrInternalError  ErrorCode = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Details string    `json:"details"`
}

type ComposeError struct {
	StatusCode int
	Code       ErrorCode
	Details    string
}

func (e *ComposeError) Error() string {
	return e.Details
}
