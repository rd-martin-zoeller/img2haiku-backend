package openai

type openAiRequest struct {
	Model       string                 `json:"model"`
	Messages    []openAiRequestMessage `json:"messages"`
	MaxTokens   int                    `json:"max_tokens"`
	Temperature float32                `json:"temperature"`
}

type openAiRequestMessage struct {
	Role    string                        `json:"role"`
	Content []openAiRequestMessageContent `json:"content"`
}

type openAiRequestMessageContent struct {
	Type     string               `json:"type"`
	Text     string               `json:"text,omitempty"`
	ImageURL *openAiImageURLField `json:"image_url,omitempty"`
}

type openAiImageURLField struct {
	URL string `json:"url"`
}

func buildRequest(prompt, base64Image string) *openAiRequest {
	return &openAiRequest{
		Model: "gpt-4o-2024-08-06",
		Messages: []openAiRequestMessage{
			{
				Role: "user",
				Content: []openAiRequestMessageContent{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: &openAiImageURLField{
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
