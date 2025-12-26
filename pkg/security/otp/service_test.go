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

		mockStorage.EXPECT().Put(ctx, mock.Anything, mock.AnythingOfType("time.Time")).
			Return(nil)

		otp, ok := manager.GenerateOTP(ctx, length, kind)

		assert.True(t, ok)
		assert.NotNil(t, otp)
		assert.Equal(t, hashedCode, otp.HashedCode)
		assert.True(t, otp.ExpiresAt.Before(time.Now().UTC().Add(time.Minute)))
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

		_, ok := manager.GenerateOTP(ctx, length, kind)

		assert.False(t, ok)
	})
}

func TestDefaultManager_Retrieve(t *testing.T) {
	ctx := context.Background()
	kind := KindEmployeePassword
	code := "existing-code"
	hashedCode := "hashed-code"

	t.Run("successfully retrieves OTP", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}
		expirationTime := time.Now().UTC().Add(time.Minute)

		expectedOTP := OneTimePassword{
			Kind:       kind,
			Code:       code,
			HashedCode: hashedCode,
			ExpiresAt:  expirationTime,
		}

		mockHasher.EXPECT().Hash(code).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Get(ctx, kind, hashedCode).Return(expectedOTP, nil)

		otp, ok := manager.Retrieve(ctx, kind, code)

		assert.True(t, ok)
		assert.Equal(t, expectedOTP, otp)
	})

	t.Run("Retrieves OTP when expired", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}
		expirationTime := time.Now().UTC().Add(-time.Minute)

		expectedOTP := OneTimePassword{
			Kind:       kind,
			Code:       code,
			HashedCode: hashedCode,
			ExpiresAt:  expirationTime,
		}

		mockHasher.EXPECT().Hash(code).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Get(ctx, kind, hashedCode).Return(expectedOTP, nil)

		_, ok := manager.Retrieve(ctx, kind, code)

		assert.False(t, ok)
	})

	t.Run("returns error when storage fails", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		mockHasher.EXPECT().Hash(code).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Get(ctx, kind, hashedCode).
			Return(OneTimePassword{}, errors.New("not found"))

		_, ok := manager.Retrieve(ctx, kind, code)

		assert.False(t, ok)
	})
}

func TestDefaultManager_Remove(t *testing.T) {
	ctx := context.Background()
	code := "existing-code"
	hashedCode := "hashed-code"

	t.Run("operation is successful", func(t *testing.T) {
		mockStorage := NewMockStorage(t)
		mockHasher := hasher.NewMockHasher(t)
		manager := &DefaultManager{
			storage:            mockStorage,
			hashingAlgo:        mockHasher,
			expirationDuration: time.Minute,
		}

		mockHasher.EXPECT().Hash(code).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Remove(ctx, hashedCode).Return(nil)

		ok := manager.Remove(ctx, code)

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

		// Expectations
		mockHasher.EXPECT().Hash(code).
			Return(hashedCode, nil).Times(1)

		mockStorage.EXPECT().Remove(ctx, hashedCode).Return(errors.New("error"))

		ok := manager.Remove(ctx, code)

		assert.False(t, ok)
	})
}
