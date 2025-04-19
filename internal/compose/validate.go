package compose

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/utils"
)

func validateRequest(r *http.Request) (types.ComposeRequest, error) {
	var req types.ComposeRequest

	if apiKey := r.Header.Get("X-API-Key"); apiKey == "" || apiKey != os.Getenv("API_KEY") {
		return req, utils.NewErr(http.StatusUnauthorized, types.ErrInternalError, "%s", "Invalid API key")
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
