package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/utils"
)

type OpenAiClient struct {
	ApiKey string
	Client *http.Client
}

const apiURL = "https://api.openai.com/v1/chat/completions"

func (c *OpenAiClient) Call(ctx context.Context, prompt, base64Image string) (types.Haiku, *types.ComposeError) {
	var haiku types.Haiku
	reqObj := buildRequest(prompt, base64Image)

	bodyBytes, err := json.Marshal(reqObj)
	if err != nil {
		return haiku, utils.NewInternalErr("Failed to encode request body: %s", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return haiku, utils.NewInternalErr("Failed to create request: %s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return haiku, utils.NewInternalErr("Failed to call OpenAI API: %s", err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return haiku, utils.NewInternalErr("OpenAI API returned an error: %s", resp.Status)
	}

	return handleResponseBody(resp)
}
