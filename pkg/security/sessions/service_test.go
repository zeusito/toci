package sessions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

func TestNewSessionFailedToCreateToken(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}

	ctx := context.Background()
	claims := PrincipalClaims{PrincipalID: "aud_id", Roles: []string{"role1", "role2"}, IsAuthenticated: true}
	expirationAsDuration := time.Hour

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
		Return("", errors.New("failed to hash")).Once()

	// Execute
	token, ok := service.NewSession(ctx, claims, expirationAsDuration)

	assert.False(t, ok)
	assert.Empty(t, token)
}

func TestNewSessionFailedToPersist(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}

	ctx := context.Background()
	claims := PrincipalClaims{PrincipalID: "aud_id", Roles: []string{"role1", "role2"}, IsAuthenticated: true}
	expirationAsDuration := time.Hour

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
		Return("hashed_token", nil).Once()

	mockStorage.EXPECT().Set(ctx, "hashed_token", mock.Anything).
		Return(errors.New("failed to persist")).Once()

	// Execute
	token, ok := service.NewSession(ctx, claims, expirationAsDuration)

	assert.False(t, ok)
	assert.Empty(t, token)
}

func TestNewSession(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}

	ctx := context.Background()
	claims := PrincipalClaims{PrincipalID: "aud_id", Roles: []string{"role1", "role2"}, IsAuthenticated: true}
	expirationAsDuration := time.Hour

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).Return("hashed_token", nil).Once()

	mockStorage.EXPECT().Set(ctx, mock.AnythingOfType("string"), mock.Anything).
		Return(nil).Once()

	// Execute
	token, ok := service.NewSession(ctx, claims, expirationAsDuration)

	assert.True(t, ok)
	assert.NotEmpty(t, token)
}

func TestGetSessionFailedToHashToken(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}
	ctx := context.Background()

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
		Return("", errors.New("failed to hash")).Once()

	// Execute
	claims := service.GetSession(ctx, "token")

	assert.False(t, claims.IsAuthenticated)
}

func TestGetSessionFailedToRetrieve(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}
	ctx := context.Background()

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).Return("hashed_token", nil).Once()

	mockStorage.EXPECT().Get(ctx, mock.AnythingOfType("string")).
		Return(nil, errors.New("failed to retrieve")).Once()

	// Execute
	claims := service.GetSession(ctx, "token")

	assert.False(t, claims.IsAuthenticated)
}

func TestGetSession(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}
	ctx := context.Background()

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).Return("hashed_token", nil).Once()

	mockStorage.EXPECT().Get(ctx, mock.AnythingOfType("string")).
		Return(&sessionData{PrincipalID: "aud_id", Roles: []string{"role1", "role2"}}, nil).Once()

	// Execute
	claims := service.GetSession(ctx, "token")

	assert.True(t, claims.IsAuthenticated)
}

func TestDeleteSessionFailedToHashToken(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}
	ctx := context.Background()

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
		Return("", errors.New("failed to hash")).Once()

	// Execute
	ok := service.RemoveSession(ctx, "token")

	assert.False(t, ok)
}

func TestDeleteSessionFailedToRetrieve(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}
	ctx := context.Background()

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).Return("hashed_token", nil).Once()

	mockStorage.EXPECT().Remove(ctx, mock.AnythingOfType("string")).
		Return(errors.New("failed to remove")).Once()

	// Execute
	ok := service.RemoveSession(ctx, "token")

	assert.False(t, ok)
}

func TestDeleteSession(t *testing.T) {
	mockStorage := NewMockStorage(t)
	mockHasher := hasher.NewMockHasher(t)
	service := &DefaultManager{
		storage:     mockStorage,
		tokenHasher: mockHasher,
	}
	ctx := context.Background()

	// Expectations
	mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).Return("hashed_token", nil).Once()

	mockStorage.EXPECT().Remove(ctx, mock.AnythingOfType("string")).
		Return(nil).Once()

	// Execute
	ok := service.RemoveSession(ctx, "token")

	assert.True(t, ok)
}
