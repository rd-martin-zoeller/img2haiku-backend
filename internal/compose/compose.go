package compose

import (
	"encoding/json"
	"net/http"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func ComposeHaiku(client types.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		composeHaiku(client, w, r)
	}
}

func composeHaiku(client types.Client, w http.ResponseWriter, r *http.Request) {
	req, err := validateRequest(r)
	if err != nil {
		writeError(w, *err)
		return
	}

	prompt, composeErr := makePrompt(req.Language, req.Tags)
	if composeErr != nil {
		writeError(w, *composeErr)
		return
	}

	haiku, err := client.Call(r.Context(), prompt, req.Base64Image)
	if err != nil {
		writeError(w, *err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(haiku)
}

func writeError(w http.ResponseWriter, error types.ComposeError) {
	w.WriteHeader(error.StatusCode)
	errorResponse := types.ErrorResponse{
		Code:    error.Code,
		Details: error.Details,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
