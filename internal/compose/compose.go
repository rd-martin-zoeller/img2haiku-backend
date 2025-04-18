package compose

import (
	"encoding/json"
	"net/http"
)

func ComposeHaiku(w http.ResponseWriter, r *http.Request) {
	_, composeErr := validateRequest(r)

	w.Header().Set("Content-Type", "application/json")

	if composeErr != nil {
		writeError(w, *composeErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	haiku := Haiku{
		Haiku: `Reges Treiben in kühler Sommernacht
		Holz trägt Glückseligkeit
		Der Gott des Donners ist im Tal`,
		Description: "Das ist ein Beispiel-Haiku",
	}
	json.NewEncoder(w).Encode(haiku)
}

func writeError(w http.ResponseWriter, error ComposeError) {
	w.WriteHeader(error.StatusCode)
	errorResponse := ErrorResponse{
		Code:    error.Code,
		Details: error.Details,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
