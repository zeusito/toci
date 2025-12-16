package sessions

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

func TestNewSession(t *testing.T) {
	mockStorage := NewMockStorage(t)
	shaHasher, _ := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	service := NewDefaultService(mockStorage, shaHasher, "test")

	ctx := context.Background()
	claims := &AuthClaims{Principal: "aud_id", Roles: []string{"role1", "role2"}, IsAuthenticated: true}
	expirationAsDuration := time.Hour

	// Expectations
	mockStorage.EXPECT().Set(ctx, mock.AnythingOfType("string"), claims, mock.AnythingOfType("time.Time")).Return(nil).Once()

	// Execute
	token, err := service.NewSession(ctx, claims, expirationAsDuration)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, strings.HasPrefix(token, "test_"), "Token should have the correct prefix")
}

func TestGetSession(t *testing.T) {
	mockStorage := NewMockStorage(t)
	shaHasher, _ := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	service := NewDefaultService(mockStorage, shaHasher, "test")

	ctx := context.Background()
	claims := &AuthClaims{Principal: "aud_id", Roles: []string{"role1", "role2"}, IsAuthenticated: true}
	opaqueToken := NewOpaqueToken("test", shaHasher)
	hashedToken, _ := opaqueToken.SecureHashedString()

	// Expectations
	mockStorage.EXPECT().Get(ctx, hashedToken).Return(claims, nil).Once()

	// Execute
	resp := service.GetSession(ctx, opaqueToken.String())

	assert.True(t, resp.IsAuthenticated)
}

func TestGetSessionInvalidToken(t *testing.T) {
	mockStorage := NewMockStorage(t)
	shaHasher, _ := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	service := NewDefaultService(mockStorage, shaHasher, "test")

	ctx := context.Background()
	opaqueToken := NewOpaqueToken("test", shaHasher)
	hashedToken, _ := opaqueToken.SecureHashedString()

	// Expectations
	mockStorage.EXPECT().Get(ctx, hashedToken).Return(nil, errors.New("record not found")).Once()

	// Execute
	resp := service.GetSession(ctx, opaqueToken.String())

	assert.False(t, resp.IsAuthenticated)
}

func TestRemoveSession(t *testing.T) {
	mockStorage := NewMockStorage(t)
	shaHasher, _ := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	service := NewDefaultService(mockStorage, shaHasher, "test")

	ctx := context.Background()
	opaqueToken := NewOpaqueToken("test", shaHasher)
	hashedToken, _ := opaqueToken.SecureHashedString()

	// Expectations
	mockStorage.EXPECT().Remove(ctx, hashedToken).Return(nil).Once()

	// Execute
	err := service.RemoveSession(ctx, opaqueToken.String())

	assert.NoError(t, err)
}
