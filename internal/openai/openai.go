package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func (c *OpenAiClient) Call(ctx context.Context, prompt, base64Image string) (types.Haiku, *types.ComposeError) {
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
	req = req.WithContext(ctx)

	if reqErr != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to create request: " + reqErr.Error(),
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := c.Client.Do(req)

	if apiErr != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to call OpenAI API: " + apiErr.Error(),
		}
	}

	defer resp.Body.Close()

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

	var openAiResponse OpenAiResponseBody
	if err := json.NewDecoder(bytes.NewBuffer(respBytes)).Decode(&openAiResponse); err != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to decode response body: " + err.Error(),
		}
	}

	if len(openAiResponse.Choices) == 0 {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "No choices found in response",
		}
	}

	answer := openAiResponse.Choices[0].Message.Content

	var haikuResponse OpenAiHaikuResponse
	if err := json.Unmarshal([]byte(answer), &haikuResponse); err != nil {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusInternalServerError,
			Code:       types.ErrInternalError,
			Details:    "Failed to unmarshal answer JSON: " + err.Error() + "\n" + answer,
		}
	}

	if haikuResponse.Error != "" {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusBadRequest,
			Code:       types.ErrInvalidRequest,
			Details:    haikuResponse.Error,
		}
	}

	if haikuResponse.Haiku == "" || haikuResponse.Description == "" {
		return haiku, &types.ComposeError{
			StatusCode: http.StatusBadRequest,
			Code:       types.ErrInvalidRequest,
			Details:    "Invalid response format: haiku or description not found " + answer,
		}
	}

	haiku.Haiku = sanitizeHaiku(haikuResponse.Haiku)
	haiku.Description = haikuResponse.Description
	return haiku, nil
}

func sanitizeHaiku(haiku string) string {
	// Sometimes, ChatGPT escapes newline characters (\n) as \\n.
	// This function replaces them with actual newlines.
	return strings.ReplaceAll(haiku, "\\n", "\n")
}
