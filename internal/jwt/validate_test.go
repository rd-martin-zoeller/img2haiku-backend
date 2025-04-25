package jwt

import (
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	keyPair1, err := GenKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	keyPair2, err := GenKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	cases := []struct {
		name       string
		privateKey string
		publicKey  string
		sub        string
		aud        string
		exp        time.Duration
		wantValid  bool
		wantError  string
	}{
		{
			name:       "valid JWT",
			privateKey: keyPair1.Private,
			publicKey:  keyPair1.Public,
			sub:        sub,
			aud:        aud,
			exp:        ttl,
			wantValid:  true,
		},
		{
			name:       "invalid signature",
			privateKey: keyPair2.Private,
			publicKey:  keyPair1.Public,
			sub:        sub,
			aud:        aud,
			exp:        ttl,
			wantValid:  false,
			wantError:  "crypto/rsa: verification error",
		},
		{
			name:       "invalid audience",
			privateKey: keyPair1.Private,
			publicKey:  keyPair1.Public,
			sub:        sub,
			aud:        "invalid-audience",
			exp:        ttl,
			wantValid:  false,
			wantError:  "invalid audience",
		},
		{
			name:       "expired token",
			privateKey: keyPair1.Private,
			publicKey:  keyPair1.Public,
			sub:        sub,
			aud:        aud,
			exp:        -time.Minute,
			wantValid:  false,
			wantError:  "Token is expired",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt, err := JWTForTesting(JWTConfig{
				KeyPair: KeyPair{
					Private: c.privateKey,
					Public:  c.publicKey,
				},
				Sub: c.sub,
				Aud: c.aud,
				Exp: c.exp,
			})
			if err != nil {
				t.Fatalf("Failed to generate JWT: %v", err)
			}

			valid, err := Validate(jwt, c.publicKey)
			if valid != c.wantValid {
				t.Fatalf("Expected JWT validity to be %v, got %v", c.wantValid, valid)
			}

			if err != nil {
				if c.wantError == "" {
					t.Fatalf("Unexpected error: %v", err)
				}
				if err.Error() != c.wantError {
					t.Fatalf("Expected error message to be %q, got %q", c.wantError, err.Error())
				}
			}
		})
	}
}
