package types

import "context"

type ComposeRequest struct {
	Language    string   `json:"language"`
	Tags        []string `json:"tags"`
	Base64Image string   `json:"base64Image"`
}

type Haiku struct {
	Haiku       string `json:"haiku"`
	Description string `json:"description"`
}

type Client interface {
	Call(ctx context.Context, prompt, base64Image string) (Haiku, *ComposeError)
}
