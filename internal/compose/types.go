package compose

type ComposeRequest struct {
	Language    string   `json:"language"`
	Tags        []string `json:"tags"`
	Base64Image string   `json:"base64Image"`
}

type Haiku struct {
	Haiku       string `json:"haiku"`
	Description string `json:"description"`
}
