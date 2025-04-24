package compose

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func ComposeHaiku(client types.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		composeHaiku(client, w, r)
	}
}

func composeHaiku(client types.Client, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := validateRequest(r)
	if err != nil {
		writeError(w, err)
		return
	}

	prompt, err := makePrompt(req.Language, req.Tags)
	if err != nil {
		writeError(w, err)
		return
	}

	haiku, err := client.Call(r.Context(), prompt, req.Base64Image)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(haiku)
}

func writeError(w http.ResponseWriter, err error) {
	var composeErr *types.ComposeError
	if errors.As(err, &composeErr) {
		w.WriteHeader(composeErr.StatusCode)
		errorResponse := types.ErrorResponse{
			Code:    composeErr.Code,
			Details: composeErr.Details,
		}
		json.NewEncoder(w).Encode(errorResponse)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := types.ErrorResponse{
			Code:    types.ErrInternalError,
			Details: "An unexpected error occurred: " + err.Error(),
		}
		json.NewEncoder(w).Encode(errorResponse)
	}
}
