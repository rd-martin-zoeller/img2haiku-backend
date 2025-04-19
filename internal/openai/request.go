package openai

func buildRequest(prompt, base64Image string) *openAiRequest {
	return &openAiRequest{
		Model: "gpt-4o-2024-08-06",
		Messages: []openAiMessage{
			{
				Role: "user",
				Content: []openAiMessageContent{
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
