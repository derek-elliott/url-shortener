package api

import (
	"crypto/rand"
	b64 "encoding/base64"
)

// GenerateToken generates a cryptographically secure random byte array of length len and encodes it into a URL-safe base 64 string
func GenerateToken(len int) (string, error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return b64.URLEncoding.EncodeToString(b), nil
}
