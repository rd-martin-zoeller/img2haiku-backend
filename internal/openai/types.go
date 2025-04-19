package openai

type OpenAiClient struct{}

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
