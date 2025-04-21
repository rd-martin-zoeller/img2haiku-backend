package jwt

import (
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
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
			privateKey: validPrivateKey,
			publicKey:  validPublicKey,
			sub:        sub,
			aud:        aud,
			exp:        ttl,
			wantValid:  true,
		},
		{
			name:       "invalid signature",
			privateKey: invalidPrivateKey,
			publicKey:  validPublicKey,
			sub:        sub,
			aud:        aud,
			exp:        ttl,
			wantValid:  false,
			wantError:  "crypto/rsa: verification error",
		},
		{
			name:       "invalid subject",
			privateKey: validPrivateKey,
			publicKey:  validPublicKey,
			sub:        "invalid-subject",
			aud:        aud,
			exp:        ttl,
			wantValid:  false,
			wantError:  "invalid subject",
		},
		{
			name:       "invalid audience",
			privateKey: validPrivateKey,
			publicKey:  validPublicKey,
			sub:        sub,
			aud:        "invalid-audience",
			exp:        ttl,
			wantValid:  false,
			wantError:  "invalid audience",
		},
		{
			name:       "expired token",
			privateKey: validPrivateKey,
			publicKey:  validPublicKey,
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

const validPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDTKXCGHegjXw8S
pL6cXnNZ8fhSrhIC7qPdO9PLd2wjC97qIhczpsJNQnv3j2Z4DQrtq3bUC4npS3VL
M3dwwN5M5jzJGZifUDlyRi6x4Kt7+5IVAB6pCxiYSPoPlMfr68ajVpg53/kypNYL
wC9vV+OrhhmAzQ2A2Omf/pnGndquARCG+Fi5JGhU1oBx66jI3CErdWt6DB1SBMzy
HSjTW4cXSqp3AP+KrfffRxzPYH9Tn9sO00qWIImGKz7vkRZ6cW02liL2Dtom9VGl
nKf+lTe/kErhgoknZ78DcBianFcGplslTQjPQCEiS/0GGSh46yFCa+tDGMmpOgG7
SjjMej9XAgMBAAECggEABMjNey4Mfb5GvECBIYNFIzZ6zspeSRrsGcEKkVnq2Gs8
vHRT4nOdHZzSHzyGColNHcmpRX1F1vd2BTBGcAatFfP16gTxv5/uaHrtNIWN84V/
L7xsqJ4PC00Rl8gvGwoHQNmSiGX21cfRu7sGhOpqYHWg3uZTSUyUFCNAi7U1TutX
uC/bmSGgOJzxeNZQxkzKRokPIHF2wI+aQ6JRigQaxK7b5WZiALqp2iiqqrrUqbyc
J7COGAPT0k7LnZS/7/TbqBTYKaVo8wMGXJRNmum/6BJNsk0F6PVcxIV0fYJ4MnKg
+ajDeV99IYPsuQl2OjXEUa4/wjRI1AwOlZw8ya3ZvQKBgQDpeWdzFWUBrORctUZx
5f+VsAtPBxic91I1CACw590UwFpcMaxzE4c4bGnZoCac2P87tSPqZvXDKojKLii5
Tc3FTPB+dKLPz0do6mLQmvXNLm0cZLAd4APhgVGWZk74gM0gBvbN+Cx1TAQScWXM
02jUOakh0PkNaN2fRRbDrJ63DQKBgQDniPDgOmN3lEfpiYwpjouxRKu74tvq5QGn
GMiQnoJCafeaGbLr0asWRD5WcUGtsfBO4+pXEcBp5IB0T7J5+1DZUPwYPaNh3gls
vS4e5hZV/ud3jQY9xJN6bnqI7x/I8ZhwL5rXFPnCLf1UHWalKeXKCOvqUE73yy6A
wiha6pz28wKBgFusXM9WVjvLLDuuvgNZAPtAjaAxNBvmDLRf+Q19bVSJlrFem8zv
nQetof5eoOqzVbyXCowuc093sxBYAYuJHkPbSw8MMyWPyQVMCxLH1b4D/bnJW1HP
tRZllaiNcXKn+GMb+Oq1CJfiCjNHrWY4mI/EOEHb8P6v711rXl3kuMk9AoGBAI6w
mWpG8afvTUZCy4uM2tBbts6q58dibNtS7cAav8I4Viy1K8wjQiIN2rEhSU3HfIbR
9UjFmuRnuzZzK1X7qP7U5xf1XKxiz0IhcLwAJsHGv1WxJqiIbi8kyQV9AQSwx7ZT
0EQ/HBEskJP3LpwZLxGM3/9ekNwrbrRRc9dcAXI7AoGBAJNluf+FlpdyOY4obZgc
8O1VUGAmX8Y5VYYBoWrpptR9lOJvAIqgg3eQNfb+jChBQc99odcHWkp7I+jjM9BI
8JsfbCEs1K+/6SzRRU1Fk7K7+Yr1q96W4wjG/AGO620mmWpBp6X1vQNyDD3j6umR
FloM7ZDtjUpRSNG70Lf3o30w
-----END PRIVATE KEY-----
`

const validPublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0ylwhh3oI18PEqS+nF5z
WfH4Uq4SAu6j3TvTy3dsIwve6iIXM6bCTUJ7949meA0K7at21AuJ6Ut1SzN3cMDe
TOY8yRmYn1A5ckYuseCre/uSFQAeqQsYmEj6D5TH6+vGo1aYOd/5MqTWC8Avb1fj
q4YZgM0NgNjpn/6Zxp3argEQhvhYuSRoVNaAceuoyNwhK3VregwdUgTM8h0o01uH
F0qqdwD/iq3330ccz2B/U5/bDtNKliCJhis+75EWenFtNpYi9g7aJvVRpZyn/pU3
v5BK4YKJJ2e/A3AYmpxXBqZbJU0Iz0AhIkv9BhkoeOshQmvrQxjJqToBu0o4zHo/
VwIDAQAB
-----END PUBLIC KEY-----
`

const invalidPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDdXHUFTikOEp9k
z0K0A3f0B8tC/foH6jC0BCMW4zsJCYl3FsDa5ZKUoOu6V4iN3fCaNpnX8QHMQYWI
vKfb2oZnwySBIvCC32YSzuLAPTcTcBBCYZqFbfNZwx6rS8ec4mAgy0/C2jGqkweA
7zmvFbzpfDMhRrjJdH9HoT60zlrJBWhDIedwaXNzE35YiZCGeyPmb2rdX96foIwL
lZtYHLmB9SG0sP8uKkdGemWrIzexEyWFj9UbG9DhVJFLFIiXfJBPpR5t9Oou93ux
IPVNniAOFQ1uhJA51Ch1kM0NY26XBX1RNQJJLHTqAdkj20Oe/w58w2GngSeZ7Hv+
d8vthXQ/AgMBAAECggEAB6ADMesW0/LFRdIz4IKME75e/JBGGBazlcfcs5GhO3b2
IsGIZCHrUi5W4GTagdSG0LEXzI3zO2d4Y5ToDVUyMwnQTJh5A3ERkY1J10hkiMlf
7gFxsq3uZ4WmnUzvc9KCcC7AsRwWAOOuqvzSllrf1oUeN8O5YsseBUgjIlRHYUw7
eWs/Y9PU7XEJGmzclZdBdmUWKnEBP8LU1qYlUuFSu12bPJ73Q2WadCcRK7ziE1kK
wNEiV1OXLVhxKD1CFTnAVTRTM0vMjpMW4TH7WP4akx2nEj4vSYLVSm7eLbuvPxGy
+iIm5bSt0sAcFq5l3iI9rSYI3/R8WH+IgeglEmVXvQKBgQD28k4WCrUAhAToG6kt
lzyNFavtDn1jYp0shiAz5sTN5lDFUDChcqTcdrIua5tTSU3DyHxdM+g1kaYTfV/O
iWsQnA/lz8ON87q0oaND180sa2DfuxngHZsgcOGMP1FMM+J/9Kbx4+dPCHlRbQ2d
Gb2P85X5etRB13y6Vx4CEcavFQKBgQDlegXn14h6Rhhyk/i1qkJEwPBVjEy3I9RN
ChTY5XK6B+aEn6xSWfSt6lC2ZDWkTuxSuqaWL3HrYPYo5zsDJJjVuGJtyeiTBI/n
oXPED4j97hbCJys8wq4VFfXGb2XE3ED6i6htDzDnH0l/hCv+WFZTHcOLq3jenf4E
Dltvat2LAwKBgB+eIH9T+Z7KSHKLcBrFPVx3BN9CNq2t55/WwHLEvjf6oCbTQJa7
Pf54OBIXdviv7wP9PGcWiUmqj0/5gnXIRGwI/0QWWNxo82PDOksqazuft/SNWR/H
yp/ZtBcn2DngfsSRR3q7ClelJxtU0iRmMk4nCvG5V0ni1DZrhw0Ox5iNAoGAHvh6
BZFMRRxivkwEPBhvezIC/bjCvdDjHUaoC6Hj+wGH9gxKyI6FfFdsb0FVEAjq0juI
sipTGK5sapbSmxj8W5PYDPM8JWNvPJbItgRWu9a/UZLRvhCUSBo/onl0Zb5IMshY
geeT9Q1+8OvYuCoZ9HvG4XnSBVGTb960LnRg1BsCgYEA5o8mDG8cZC/NrRsRPzV3
e3gEWavCi0Yr1nuflQYHIGxp87RFLPgM1DO2vqpXgns1LzyWV2rIhYGGWR53f348
17AxvYFFmDU+YTCVOmmbEWvDRrM6FqOYgqT07TTwZy1gPkO7diicNb/oI7sYL/3y
xSA2kIcp5Ar7vi8fRehcwJE=
-----END PRIVATE KEY-----
`
