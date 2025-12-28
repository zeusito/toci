package otp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

func TestDefaultManager_GenerateOTP(t *testing.T) {
	ctx := context.Background()
	kind := KindPanelistPassword
	length := 6

	t.Run("successfully generates and stores OTP", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		hashedCode := "hashed-code"

		// Expectations
		mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Put(ctx, kind, "john@example.com",
			mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

		code, ok := manager.GenerateCode(ctx, length, kind, "john@example.com")

		assert.True(t, ok)
		assert.NotEmpty(t, code)
	})

	t.Run("returns error if hashing fails", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		// Expectations
		mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
			Return("", errors.New("hashing error")).Times(1)

		code, ok := manager.GenerateCode(ctx, length, KindEmployeePassword, "john@example.com")

		assert.False(t, ok)
		assert.Empty(t, code)
	})
}

func TestDefaultManager_Validate(t *testing.T) {
	ctx := context.Background()
	kind := KindEmployeePassword
	principal := "john@example.com"
	code := "existing-code"
	hashedCode := "hashed-code"

	t.Run("successfully validates OTP", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		// Expectations
		mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Get(ctx, kind, principal).
			Return(&otpData{
				ID:        hashedCode,
				Kind:      kind,
				Principal: principal,
				ExpiresAt: time.Now().UTC().Add(time.Minute),
			}, nil).Times(1)

		ok := manager.VerifyCode(ctx, kind, principal, code)

		assert.True(t, ok)
	})

	t.Run("validation fails when hashing fails", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		// Expectations
		mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
			Return("", errors.New("hashing error")).Times(1)

		ok := manager.VerifyCode(ctx, kind, principal, code)

		assert.False(t, ok)
	})

	t.Run("validation fails when no record is found", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		// Expectations
		mockHasher.EXPECT().Hash(mock.AnythingOfType("string")).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Get(ctx, kind, principal).
			Return(nil, errors.New("record not found")).Times(1)

		ok := manager.VerifyCode(ctx, kind, principal, code)

		assert.False(t, ok)
	})
}

func TestDefaultManager_Remove(t *testing.T) {
	ctx := context.Background()
	kind := KindEmployeePassword
	principal := "john@example.com"

	t.Run("operation is successful", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		// Expectations
		mockStorage.EXPECT().Remove(ctx, kind, principal).Return(nil)

		ok := manager.Remove(ctx, kind, principal)

		assert.True(t, ok)
	})

	t.Run("returns error when function fails", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		mockStorage.EXPECT().Remove(ctx, kind, principal).Return(errors.New("error"))

		ok := manager.Remove(ctx, kind, principal)

		assert.False(t, ok)
	})
}
