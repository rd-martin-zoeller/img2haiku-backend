package function

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/compose"
)

func init() {
	functions.HTTP("ComposeHaiku", compose.ComposeHaiku)
}
