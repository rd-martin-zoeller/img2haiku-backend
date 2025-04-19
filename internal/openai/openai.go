package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/utils"
)

func (c *OpenAiClient) Call(ctx context.Context, prompt, base64Image string) (types.Haiku, *types.ComposeError) {
	var haiku types.Haiku
	reqObj := makeRequestObject(prompt, base64Image)

	bodyBytes, err := json.Marshal(reqObj)
	if err != nil {
		return haiku, utils.NewInternalErr("Failed to encode request body: %s", err.Error())
	}

	req, reqErr := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	req = req.WithContext(ctx)

	if reqErr != nil {
		return haiku, utils.NewInternalErr("Failed to create request: %s", reqErr.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := c.Client.Do(req)

	if apiErr != nil {
		return haiku, utils.NewInternalErr("Failed to call OpenAI API: %s", apiErr.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return haiku, utils.NewInternalErr("OpenAI API returned an error: %s", resp.Status)
	}

	return handleResponseBody(resp)
}

func makeRequestObject(prompt, base64Image string) *OpenAiRequest {
	return &OpenAiRequest{
		Model: "gpt-4o-2024-08-06",
		Messages: []OpenAiMessage{
			{
				Role: "user",
				Content: []OpenAiMessageContent{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: &OpenAiImageURLField{
							URL: "data:image/jpeg;base64," + base64Image,
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

	var openAiResponse OpenAiResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&openAiResponse); err != nil {
		return haiku, utils.NewInternalErr("Failed to decode response body: %s", err.Error())
	}

	if len(openAiResponse.Choices) == 0 {
		return haiku, utils.NewInternalErr("%s", "No choices found in response")
	}

	answer := openAiResponse.Choices[0].Message.Content

	var haikuResponse OpenAiHaikuResponse
	if err := json.Unmarshal([]byte(answer), &haikuResponse); err != nil {
		return haiku, utils.NewInternalErr("Failed to unmarshal answer JSON: %s\n%s", err.Error(), answer)
	}

	if haikuResponse.Error != "" {
		return haiku, utils.NewErr(http.StatusBadRequest, types.ErrInvalidRequest, "%s", haikuResponse.Error)
	}

	if haikuResponse.Haiku == "" || haikuResponse.Description == "" {
		return haiku, utils.NewErr(http.StatusBadRequest, types.ErrInvalidRequest, "Invalid response format: haiku or description not found %s", answer)
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
