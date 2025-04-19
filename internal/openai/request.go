package openai

type request struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float32       `json:"temperature"`
}

type chatMessage struct {
	Role    string           `json:"role"`
	Content []messageContent `json:"content"`
}

type messageContent struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageUrl `json:"image_url,omitempty"`
}

type imageUrl struct {
	URL string `json:"url"`
}

func buildRequest(prompt, base64Image string) *request {
	return &request{
		Model: "gpt-4o-2024-08-06",
		Messages: []chatMessage{
			{
				Role: "user",
				Content: []messageContent{
					{
						Type: "text",
						Text: prompt,
					},
					{
						Type: "image_url",
						ImageURL: &imageUrl{
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
