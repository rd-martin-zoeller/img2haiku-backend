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

func init() {
	functions.HTTP("ComposeHaiku", composeHaiku)
}

func composeHaiku(w http.ResponseWriter, r *http.Request) {
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
