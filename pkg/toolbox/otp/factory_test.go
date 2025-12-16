package otp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

//nolint:gosec
func TestNew(t *testing.T) {
	validSecret := "dGVzdC1zZWNyZXQ=" // base64 for "test-secret"
	theHasher, err := hasher.NewHmacSHA256(validSecret)
	assert.NoError(t, err)

	generator := NewOneTimePasswordGenerator(theHasher)
	token, err := generator.New(10)

	assert.NoError(t, err)
	assert.Len(t, token.Secret, 10, "Expected secret length 10")
	assert.NotEmpty(t, token.HashedSecret, "Expected non-empty hashed value")
}

//nolint:gosec
func TestVerify(t *testing.T) {
	validSecret := "dGVzdC1zZWNyZXQ=" // base64 for "test-secret"
	theHasher, err := hasher.NewHmacSHA256(validSecret)
	assert.NoError(t, err)

	generator := NewOneTimePasswordGenerator(theHasher)
	token, err := generator.New(10)

	assert.NoError(t, err)
	assert.True(t, generator.Verify(token.HashedSecret, token.Secret), "Expected verification to succeed")
	assert.False(t, generator.Verify(token.HashedSecret, "wrongsecret"), "Expected verification to fail with wrong secret")
}
