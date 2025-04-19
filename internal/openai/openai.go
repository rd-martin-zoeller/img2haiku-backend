package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

type OpenAiRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float32         `json:"temperature"`
}

type OpenAiMessage struct {
	Role    string           `json:"role"`
	Content []map[string]any `json:"content"`
}

type OpenAiClient struct{}

func (c *OpenAiClient) Call(prompt, base64Image string) (types.Haiku, *types.ComposeError) {
	var haiku types.Haiku
	reqObj := makeRequestObject(prompt, base64Image)

	bodyBytes, err := json.Marshal(reqObj)
	if err != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to encode request body: " + err.Error(),
		}
	}

	req, reqErr := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if reqErr != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to create request: " + reqErr.Error(),
		}
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := http.DefaultClient.Do(req)
	if apiErr != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to call OpenAI API: " + apiErr.Error(),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return haiku, &types.ComposeError{
			StatusCode: resp.StatusCode,
			Code:       types.ErrInternalError,
			Details:    "OpenAI API returned an error: " + resp.Status,
		}
	}

	return handleResponseBody(resp)
}

func makeRequestObject(prompt, base64Image string) *OpenAiRequest {
	return &OpenAiRequest{
		Model: "gpt-4o-2024-08-06",
		Messages: []OpenAiMessage{
			{
				Role: "user",
				Content: []map[string]any{
					{
						"type": "text",
						"text": prompt,
					},
					{
						"type": "image_url",
						"image_url": map[string]any{
							"url": "data:image/jpeg;base64," + base64Image,
						},
					},
				},
			},
		},
		MaxTokens:   200,
		Temperature: 1.0,
	}
}

func handleResponseBody(resp *http.Response) (types.Haiku, *types.ComposeError) {
	var haiku types.Haiku

	respBytes, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to read response body: " + ioErr.Error(),
		}
	}

	var respMap map[string]any
	if err := json.NewDecoder(bytes.NewBuffer(respBytes)).Decode(&respMap); err != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to decode response body: " + err.Error(),
		}
	}

	rawChoices, ok := respMap["choices"]
	if !ok {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "No choices found in response",
		}
	}

	choices, ok := rawChoices.([]any)
	if !ok {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "No choices found in response",
		}
	}

	if len(choices) == 0 {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "No choices found in response",
		}
	}

	answer := choices[0].(map[string]any)["message"].(map[string]any)["content"].(string)

	var answerJSON map[string]any
	if err := json.Unmarshal([]byte(answer), &answerJSON); err != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to unmarshal answer JSON: " + err.Error() + "\n" + answer,
		}
	}

	answerHaiku, haikuOk := stringField(answerJSON, "haiku")
	answerDescription, descriptionOk := stringField(answerJSON, "description")
	answerError, errorOk := stringField(answerJSON, "error")

	if errorOk && answerError != "" {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusBadRequest,
			Code:       types.ErrInvalidRequest,
			Details:    string(answerError),
		}
	}

	if !haikuOk || !descriptionOk {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusBadRequest,
			Code:       types.ErrInvalidRequest,
			Details:    "Invalid response format: haiku or description not found " + answer,
		}
	}

	haiku.Haiku = answerHaiku
	haiku.Description = answerDescription
	return haiku, nil
}

func stringField(m map[string]any, key string) (string, bool) {
	raw, ok := m[key]
	if !ok || raw == nil {
		return "", false
	}
	str, ok := raw.(string)
	if ok {
		return str, true
	}

	return fmt.Sprintf("%v", raw), true
}
