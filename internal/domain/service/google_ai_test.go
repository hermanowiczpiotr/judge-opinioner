package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGoogleAiService(t *testing.T) {
	service, err := NewGoogleAiService("test", "test", "test")

	assert.NoError(t, err)

	var expectedType *GoogleAiService
	assert.IsTypef(t, expectedType, service, "")
}
