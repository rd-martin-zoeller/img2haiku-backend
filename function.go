package function

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type Haiku struct {
	Haiku       string `json:"haiku"`
	Description string `json:"description"`
}

type ErrorCode string

const (
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrInternalError  ErrorCode = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Details string    `json:"details"`
}

type ComposeRequest struct {
	Language    string   `json:"language"`
	Tags        []string `json:"tags"`
	Base64Image string   `json:"base64_image"`
}

func init() {
	functions.HTTP("ComposeHaiku", composeHaiku)
}

func composeHaiku(w http.ResponseWriter, r *http.Request) {
	_, composeErr := validateRequest(r)

	if composeErr != nil {
		writeError(w, *composeErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	haiku := Haiku{
		Haiku: `Reges Treiben in kühler Sommernacht
		Holz trägt Glückseligkeit
		Der Gott des Donners ist im Tal`,
		Description: "Das ist ein Beispiel-Haiku",
	}
	json.NewEncoder(w).Encode(haiku)
}

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

type ComposeError struct {
	StatusCode int
	Code       ErrorCode
	Details    string
}

func writeError(w http.ResponseWriter, error ComposeError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.StatusCode)
	errorResponse := ErrorResponse{
		Code:    error.Code,
		Details: error.Details,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
