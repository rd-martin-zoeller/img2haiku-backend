package function

import (
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/compose"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/openai"
)

func init() {
	functions.HTTP("ComposeHaiku", compose.ComposeHaiku(&openai.OpenAiClient{
		ApiKey: os.Getenv("OPENAI_API_KEY"),
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}))
}
