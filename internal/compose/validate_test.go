package compose

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	cases := []struct {
		name           string
		httpMethod     string
		body           *ComposeRequest
		wantStatusCode int
		wantErrorCode  ErrorCode
		wantDetails    string
	}{
		{
			name:           "method is not POST",
			httpMethod:     "GET",
			wantStatusCode: 405,
			wantErrorCode:  ErrInternalError,
			wantDetails:    "Method not allowed",
		},
		{
			name:           "body is nil",
			httpMethod:     "POST",
			wantStatusCode: 500,
			wantErrorCode:  ErrInternalError,
			wantDetails:    "Failed to decode request body: EOF",
		},
		{
			name:           "body is empty",
			httpMethod:     "POST",
			body:           &ComposeRequest{},
			wantStatusCode: 500,
			wantErrorCode:  ErrInternalError,
			wantDetails:    "Language is required",
		},
		{
			name:           "language is empty",
			httpMethod:     "POST",
			body:           &ComposeRequest{Language: ""},
			wantStatusCode: 500,
			wantErrorCode:  ErrInternalError,
			wantDetails:    "Language is required",
		},
		{
			name:           "base64 image is empty",
			httpMethod:     "POST",
			body:           &ComposeRequest{Language: "English", Base64Image: ""},
			wantStatusCode: 500,
			wantErrorCode:  ErrInternalError,
			wantDetails:    "Base64 image is required",
		},
	}

	for _, c := range cases {
		c := c // capture range variable
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			bodyBytes := requestJSONHelper(t, c.body)

			req := httptest.NewRequest(c.httpMethod, "/", strings.NewReader(string(bodyBytes)))

			_, err := validateRequest(req)

			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			if err.StatusCode != c.wantStatusCode {
				t.Errorf("Expected status code %d, got %d", c.wantStatusCode, err.StatusCode)
			}

			if err.Code != c.wantErrorCode {
				t.Errorf("Expected error code %s, got %s", c.wantErrorCode, err.Code)
			}

			if err.Details != c.wantDetails {
				t.Errorf("Expected error details %s, got %s", c.wantDetails, err.Details)
			}
		})
	}
}

func requestJSONHelper(t *testing.T, body *ComposeRequest) []byte {
	t.Helper()

	var bodyBytes []byte
	if body != nil {
		var marshalErr error
		bodyBytes, marshalErr = json.Marshal(body)
		if marshalErr != nil {
			t.Fatalf("failed to marshal ComposeRequest: %v", marshalErr)
		}
	}

	return bodyBytes
}
