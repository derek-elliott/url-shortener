package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := generateToken(6)
	assert.Nil(t, err)
	assert.NotEqual(t, "", token, "Should not be an empty string")
}
