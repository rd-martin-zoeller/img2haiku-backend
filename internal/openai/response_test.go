package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestHandleResponseBody(t *testing.T) {
	cases := []struct {
		name             string
		responseBody     openAiResponseBody
		wantHaiku        string
		wantErrorMessage string
	}{
		{
			name:             "no choices field",
			responseBody:     openAiResponseBody{},
			wantErrorMessage: "No choices found in response",
		},
		{
			name:             "no choices",
			responseBody:     openAiResponseBody{Choices: []openAiResponseChoice{}},
			wantErrorMessage: "No choices found in response",
		},
		{
			name: "not a JSON content",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: "EXAMPLE_HAIKU",
						},
					},
				},
			},
			wantErrorMessage: "Failed to unmarshal answer JSON: invalid character 'E' looking for beginning of value\nEXAMPLE_HAIKU",
		},
		{
			name: "error JSON",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: `{"error":"EXAMPLE_ERROR"}`,
						},
					},
				},
			},
			wantErrorMessage: "EXAMPLE_ERROR",
		},
		{
			name: "no relevant JSON field",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: `{"some_other_field":"SOME_OTHER_FIELD"}`,
						},
					},
				},
			},
			wantErrorMessage: `Invalid response format: haiku or description not found {"some_other_field":"SOME_OTHER_FIELD"}`,
		},
		{
			name: "description missing in JSON",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: `{"haiku":"EXAMPLE_HAIKU"}`,
						},
					},
				},
			},
			wantErrorMessage: `Invalid response format: haiku or description not found {"haiku":"EXAMPLE_HAIKU"}`,
		},
		{
			name: "haiku missing in JSON",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: `{"description":"EXAMPLE_DESCRIPTION"}`,
						},
					},
				},
			},
			wantErrorMessage: `Invalid response format: haiku or description not found {"description":"EXAMPLE_DESCRIPTION"}`,
		},
		{
			name: "valid JSON",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: `{"description":"EXAMPLE_DESCRIPTION","haiku":"EXAMPLE_HAIKU"}`,
						},
					},
				},
			},
			wantHaiku: "EXAMPLE_HAIKU",
		},
		{
			name: "sanitizes valid JSON",
			responseBody: openAiResponseBody{
				Choices: []openAiResponseChoice{
					{
						Message: openAiResponseMessage{
							Content: `{"description":"EXAMPLE_DESCRIPTION","haiku":"EXAMPLE_HAIKU\\nEXAMPLE_HAIKU"}`,
						},
					},
				},
			},
			wantHaiku: "EXAMPLE_HAIKU\nEXAMPLE_HAIKU",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(c.responseBody)
			if err != nil {
				t.Fatalf("Failed to convert object to JSON: %v", err)
			}

			httpResponse := http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(bodyBytes)),
			}

			haiku, composeErr := handleResponseBody(&httpResponse)
			if c.wantErrorMessage != "" {
				if composeErr == nil {
					t.Fatalf("Expected error, got nil")
				} else {
					if composeErr.Details != c.wantErrorMessage {
						t.Errorf("Expected error details: %s\nActual error details: %s", c.wantErrorMessage, composeErr.Details)
					}
				}
			}

			if c.wantHaiku != "" {
				if haiku.Haiku != c.wantHaiku {
					t.Errorf("Expected haiku %s, got %s", c.wantHaiku, haiku.Haiku)
				}
			}
		})
	}
}
