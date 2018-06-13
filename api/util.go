package api

import (
	b64 "encoding/base64"
	"strconv"
)

// ConvertIDToToken converts an id into a URL base 64 encoded string
func ConvertIDToToken(id int) string {
	return b64.URLEncoding.EncodeToString([]byte(strconv.Itoa(id)))
}

// ConvertTokenToID converts a base 64 encoded string into an id
func ConvertTokenToID(token string) (int, error) {
	idStr, err := b64.URLEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(string(idStr))
	if err != nil {
		return 0, err
	}
	return id, nil
}
