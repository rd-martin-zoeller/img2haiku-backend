package compose

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

const validApiKey, invalidApiKey = "valid_api_key", "invalid_api_key"

func TestValidateRequest(t *testing.T) {
	cases := []struct {
		name           string
		httpMethod     string
		body           *types.ComposeRequest
		apiKey         string
		wantStatusCode int
		wantErrorCode  types.ErrorCode
		wantDetails    string
	}{
		{
			name:           "API key is missing",
			httpMethod:     "GET",
			wantStatusCode: 401,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Invalid API key",
		},
		{
			name:           "API key is invalid",
			httpMethod:     "GET",
			apiKey:         invalidApiKey,
			wantStatusCode: 401,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Invalid API key",
		},
		{
			name:           "method is not POST",
			httpMethod:     "GET",
			apiKey:         validApiKey,
			wantStatusCode: 405,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Method not allowed",
		},
		{
			name:           "body is nil",
			httpMethod:     "POST",
			apiKey:         validApiKey,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Failed to decode request body: EOF",
		},
		{
			name:           "body is empty",
			httpMethod:     "POST",
			body:           &types.ComposeRequest{},
			apiKey:         validApiKey,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Language is required",
		},
		{
			name:           "language is empty",
			httpMethod:     "POST",
			body:           &types.ComposeRequest{Language: ""},
			apiKey:         validApiKey,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Language is required",
		},
		{
			name:           "base64 image is empty",
			httpMethod:     "POST",
			body:           &types.ComposeRequest{Language: "English", Base64Image: ""},
			apiKey:         validApiKey,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Base64 image is required",
		},
	}

	for _, c := range cases {
		c := c // capture range variable
		t.Setenv("API_KEY", validApiKey)

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			bodyBytes := requestJSONHelper(t, c.body)

			req := httptest.NewRequest(c.httpMethod, "/", strings.NewReader(string(bodyBytes)))
			req.Header.Set("X-API-Key", c.apiKey)

			_, err := validateRequest(req)

			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			var composeErr *types.ComposeError
			if errors.As(err, &composeErr) {
				if composeErr.StatusCode != c.wantStatusCode {
					t.Errorf("Expected status code %d, got %d", c.wantStatusCode, composeErr.StatusCode)
				}

				if composeErr.Code != c.wantErrorCode {
					t.Errorf("Expected error code %s, got %s", c.wantErrorCode, composeErr.Code)
				}

				if composeErr.Details != c.wantDetails {
					t.Errorf("Expected error details %s, got %s", c.wantDetails, composeErr.Details)
				}
			} else {
				t.Fatalf("Expected ComposeError, got %T, %s", err, err.Error())
			}
		})
	}
}

func requestJSONHelper(t *testing.T, body *types.ComposeRequest) []byte {
	t.Helper()

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal ComposeRequest: %v", err)
		}

		return bodyBytes
	}

	return []byte{}
}
