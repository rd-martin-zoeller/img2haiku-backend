package openai

import "net/http"

type OpenAiClient struct {
	ApiKey string
	Client *http.Client
}

type openAiRequest struct {
	Model       string          `json:"model"`
	Messages    []openAiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float32         `json:"temperature"`
}

type openAiMessage struct {
	Role    string                 `json:"role"`
	Content []openAiMessageContent `json:"content"`
}

type openAiMessageContent struct {
	Type     string               `json:"type"`
	Text     string               `json:"text,omitempty"`
	ImageURL *openAiImageURLField `json:"image_url,omitempty"`
}

type openAiImageURLField struct {
	URL string `json:"url"`
}

type openAiResponseBody struct {
	Choices []openAiResponseChoice `json:"choices"`
}

type openAiResponseChoice struct {
	Message openAiResponseMessage `json:"message"`
}

type openAiResponseMessage struct {
	Content string `json:"content"`
}

type openAiHaikuResponse struct {
	Haiku       string `json:"haiku"`
	Description string `json:"description"`
	Error       string `json:"error"`
}
