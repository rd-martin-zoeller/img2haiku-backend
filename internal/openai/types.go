package openai

import "net/http"

type OpenAiClient struct {
	ApiKey string
	Client *http.Client
}

type OpenAiRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float32         `json:"temperature"`
}

type OpenAiMessage struct {
	Role    string                 `json:"role"`
	Content []OpenAiMessageContent `json:"content"`
}

type OpenAiMessageContent struct {
	Type     string               `json:"type"`
	Text     string               `json:"text,omitempty"`
	ImageURL *OpenAiImageURLField `json:"image_url,omitempty"`
}

type OpenAiImageURLField struct {
	URL string `json:"url"`
}

type OpenAiResponseBody struct {
	Choices []OpenAiResponseChoice `json:"choices"`
}

type OpenAiResponseChoice struct {
	Message OpenAiResponseMessage `json:"message"`
}

type OpenAiResponseMessage struct {
	Content string `json:"content"`
}

type OpenAiHaikuResponse struct {
	Haiku       string `json:"haiku"`
	Description string `json:"description"`
	Error       string `json:"error"`
}
