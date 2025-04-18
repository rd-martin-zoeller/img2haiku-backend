package compose

import (
	"encoding/json"
	"net/http"
)

func validateRequest(r *http.Request) (ComposeRequest, *ComposeError) {
	var req ComposeRequest

	if r.Method != http.MethodPost {
		return req, &ComposeError{
			StatusCode: http.StatusMethodNotAllowed,
			Code:       ErrInternalError,
			Details:    "Method not allowed",
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, &ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       ErrInternalError,
			Details:    "Failed to decode request body: " + err.Error(),
		}
	}

	if req.Language == "" {
		return req, &ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       ErrInternalError,
			Details:    "Language is required",
		}
	}

	if req.Base64Image == "" {
		return req, &ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       ErrInternalError,
			Details:    "Base64 image is required",
		}
	}

	return req, nil
}
