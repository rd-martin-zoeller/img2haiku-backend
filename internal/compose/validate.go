package compose

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/jwt"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/utils"
)

func validateRequest(r *http.Request) (types.ComposeRequest, error) {
	var req types.ComposeRequest

	if err := validateAuthHeader(r); err != nil {
		return req, err
	}

	if r.Method != http.MethodPost {
		return req, utils.NewErr(http.StatusMethodNotAllowed, types.ErrInternalError, "%s", "Method not allowed")
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, utils.NewInternalErr("%s", "Failed to decode request body: "+err.Error())
	}

	if req.Language == "" {
		return req, utils.NewInternalErr("%s", "Language is required")
	}

	if req.Base64Image == "" {
		return req, utils.NewInternalErr("%s", "Base64 image is required")
	}

	return req, nil
}

func validateAuthHeader(r *http.Request) error {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return utils.NewErr(http.StatusUnauthorized, types.ErrInternalError, "%s", "Authorization header is required")
	}
	if len(auth) < 7 || auth[:7] != "Bearer " {
		return utils.NewErr(http.StatusUnauthorized, types.ErrInternalError, "%s", "Authorization header must start with 'Bearer '")
	}
	token := auth[7:]
	valid, err := jwt.Validate(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return utils.NewErr(http.StatusUnauthorized, types.ErrInternalError, "%s", "Invalid JWT token: "+err.Error())
	}
	if !valid {
		return utils.NewErr(http.StatusUnauthorized, types.ErrInternalError, "%s", "Invalid JWT token")
	}

	return nil
}
