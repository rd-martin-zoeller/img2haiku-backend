package compose

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func validateRequest(r *http.Request) (types.ComposeRequest, *types.ComposeError) {
	var req types.ComposeRequest

	if apiKey := r.Header.Get("X-API-Key"); apiKey == "" || apiKey != os.Getenv("API_KEY") {
		return req, &types.ComposeError{
			StatusCode: http.StatusUnauthorized,
			Code:       types.ErrInternalError,
			Details:    "Invalid API key",
		}
	}

	if r.Method != http.MethodPost {
		return req, &types.ComposeError{
			StatusCode: http.StatusMethodNotAllowed,
			Code:       types.ErrInternalError,
			Details:    "Method not allowed",
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to decode request body: " + err.Error(),
		}
	}

	if req.Language == "" {
		return req, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Language is required",
		}
	}

	if req.Base64Image == "" {
		return req, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Base64 image is required",
		}
	}

	return req, nil
}
