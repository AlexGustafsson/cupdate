package vapid

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const AuthorizationScheme = "vapid"

func FormatAuthorizationHeader(token string, key ecdsa.PublicKey) string {
	k := base64.RawURLEncoding.EncodeToString(elliptic.MarshalCompressed(elliptic.P256(), key.X, key.Y))
	return fmt.Sprintf("vapid t=%s, k=%s", token, k)
}

// - expires MUST be less than 24 hours.
func NewToken(audience string, expires time.Time, subject string, key *ecdsa.PrivateKey) (string, error) {
	header := map[string]any{
		"typ": "JWT",
		"alg": "ES256",
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claims := map[string]any{
		"aud": audience,
		"exp": expires.Unix(),
		"sub": subject,
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	jwt := base64.RawURLEncoding.EncodeToString(headerBytes) + "." + base64.RawStdEncoding.EncodeToString(claimsBytes)

	hash := sha256.Sum256([]byte(jwt))

	signature, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		return "", err
	}

	return jwt + "." + base64.RawStdEncoding.EncodeToString(signature), nil
}
