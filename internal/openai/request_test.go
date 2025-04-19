package openai

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildRequest(t *testing.T) {
	obj := buildRequest("EXAMPLE_PROMPT", "EXAMPLE_BASE64_IMAGE")

	bodyBytes, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Failed to convert object to JSON: %v", err)
	}

	json := string(bodyBytes)

	want := `{"model":"gpt-4o-2024-08-06","messages":[{"role":"user","content":[{"type":"text","text":"EXAMPLE_PROMPT"},{"type":"image_url","image_url":{"url":"data:image/jpeg;base64,EXAMPLE_BASE64_IMAGE"}}]}],"max_tokens":200,"temperature":1}`

	if strings.Compare(json, want) != 0 {
		t.Errorf("Expected JSON: %s, got: %s", want, json)
	}
}
