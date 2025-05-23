package openai

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/utils"
)

type response struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

type message struct {
	Content string `json:"content"`
}

type haikuAnswer struct {
	Haiku       string `json:"haiku"`
	Description string `json:"description"`
	Error       string `json:"error"`
}

func handleResponseBody(resp *http.Response) (types.Haiku, error) {
	var haiku types.Haiku

	var openAiResponse response
	if err := json.NewDecoder(resp.Body).Decode(&openAiResponse); err != nil {
		return haiku, utils.NewInternalErr("Failed to decode response body: %s", err.Error())
	}

	if len(openAiResponse.Choices) == 0 {
		return haiku, utils.NewInternalErr("%s", "No choices found in response")
	}

	answer := openAiResponse.Choices[0].Message.Content

	var haikuResponse haikuAnswer
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
