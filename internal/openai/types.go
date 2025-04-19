package openai

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
