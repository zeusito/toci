package otp

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/pkg/toolbox"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type DefaultManager struct {
	hashingAlgo        hasher.Hasher
	storage            Storage
	expirationDuration time.Duration
}

func (s *DefaultManager) GenerateOTP(ctx context.Context, length int, kind Kind) (OneTimePassword, bool) {
	code := toolbox.SecureRandomString(length)
	now := time.Now().UTC()

	hashedCode, err := s.hashingAlgo.Hash(code)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash code")
		return OneTimePassword{}, false
	}

	codeHolder := OneTimePassword{
		HashedCode: hashedCode,
		Code:       code,
		Kind:       kind,
		ExpiresAt:  now.Add(s.expirationDuration),
	}

	// Persist the OTP
	err = s.storage.Put(ctx, codeHolder, codeHolder.ExpiresAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to persist OTP")
		return OneTimePassword{}, false
	}

	return codeHolder, true
}

func (s *DefaultManager) Retrieve(ctx context.Context, kind Kind, code string) (OneTimePassword, bool) {
	now := time.Now().UTC()

	hashedCode, err := s.hashingAlgo.Hash(code)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash code")
		return OneTimePassword{}, false
	}

	codeHolder, err := s.storage.Get(ctx, kind, hashedCode)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve OTP")
		return OneTimePassword{}, false
	}

	// Check if OTP is expired
	if codeHolder.ExpiresAt.Before(now) {
		log.Error().Msg("OTP is expired")
		return OneTimePassword{}, false
	}

	// Add missing code
	codeHolder.Code = code

	return codeHolder, true
}

func (s *DefaultManager) Remove(ctx context.Context, code string) bool {
	hashedCode, err := s.hashingAlgo.Hash(code)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash code")
		return false
	}

	err = s.storage.Remove(ctx, hashedCode)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove OTP")
		return false
	}

	return true
}
