package utils

import (
	"fmt"
	"net/http"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func NewErr(status int, code types.ErrorCode, msgFmt string, args ...any) error {
	return &types.ComposeError{StatusCode: status, Code: code, Details: fmt.Sprintf(msgFmt, args...)}
}

func NewInternalErr(msgFmt string, args ...any) error {
	return NewErr(http.StatusInternalServerError, types.ErrInternalError, msgFmt, args...)
}
