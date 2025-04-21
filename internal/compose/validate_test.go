package compose

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rd-martin-zoeller/img2haiku-backend/internal/jwt"
	"github.com/rd-martin-zoeller/img2haiku-backend/internal/types"
)

func TestValidateRequest(t *testing.T) {
	keyPair, err := jwt.GenKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	invalidToken := token(t, keyPair, -time.Minute)
	validToken := token(t, keyPair, time.Minute)
	t.Setenv("JWT_SECRET", keyPair.Public)
	cases := []struct {
		name           string
		httpMethod     string
		body           *types.ComposeRequest
		token          string
		wantStatusCode int
		wantErrorCode  types.ErrorCode
		wantDetails    string
	}{
		{
			name:           "token is missing",
			httpMethod:     "GET",
			wantStatusCode: 401,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Invalid JWT token: token contains an invalid number of segments",
		},
		{
			name:           "token is invalid",
			httpMethod:     "GET",
			token:          invalidToken,
			wantStatusCode: 401,
			wantErrorCode:  types.ErrAuthExpired,
			wantDetails:    "Token is expired",
		},
		{
			name:           "method is not POST",
			httpMethod:     "GET",
			token:          validToken,
			wantStatusCode: 405,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Method not allowed",
		},
		{
			name:           "body is nil",
			httpMethod:     "POST",
			token:          validToken,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Failed to decode request body: EOF",
		},
		{
			name:           "body is empty",
			httpMethod:     "POST",
			body:           &types.ComposeRequest{},
			token:          validToken,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Language is required",
		},
		{
			name:           "language is empty",
			httpMethod:     "POST",
			body:           &types.ComposeRequest{Language: ""},
			token:          validToken,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Language is required",
		},
		{
			name:           "base64 image is empty",
			httpMethod:     "POST",
			body:           &types.ComposeRequest{Language: "English", Base64Image: ""},
			token:          validToken,
			wantStatusCode: 500,
			wantErrorCode:  types.ErrInternalError,
			wantDetails:    "Base64 image is required",
		},
	}

	for _, c := range cases {
		c := c // capture range variable
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			bodyBytes := requestJSONHelper(t, c.body)

			req := httptest.NewRequest(c.httpMethod, "/", strings.NewReader(string(bodyBytes)))
			req.Header.Set("Authorization", "Bearer "+c.token)

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

func token(t *testing.T, keyPair jwt.KeyPair, exp time.Duration) string {
	t.Helper()

	token, err := jwt.JWTForTesting(jwt.JWTConfig{
		KeyPair: keyPair,
		Sub:     "img2haiku-backend-demo",
		Aud:     "img2haiku-backend",
		Exp:     exp,
	})

	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	return token
}
