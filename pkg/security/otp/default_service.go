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

// GenerateCode generates a random code of the specified length and kind
func (s *DefaultManager) GenerateCode(ctx context.Context, length int, kind Kind, principal string) (string, bool) {
	code := toolbox.SecureRandomString(length)
	now := time.Now().UTC()

	hashedCode, err := s.hashingAlgo.Hash(code)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash code")
		return "", false
	}

	// Persist the OTP
	err = s.storage.Put(ctx, kind, principal, hashedCode, now.Add(s.expirationDuration))
	if err != nil {
		log.Error().Err(err).Msg("failed to persist OTP")
		return "", false
	}

	return code, true
}

// VerifyCode verifies the code of the specified kind and principal.
// By default, only the last code from the combined kind and principal is valid.
// Expiration is checked at the storage level.
func (s *DefaultManager) VerifyCode(ctx context.Context, kind Kind, principal string, code string) bool {
	hashedCode, err := s.hashingAlgo.Hash(code)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash code")
		return false
	}

	record, err := s.storage.Get(ctx, kind, principal)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve OTP")
		return false
	}

	// check if the hashes match
	if record.ID != hashedCode {
		log.Error().Msg("hashes do not match")
		return false
	}

	return true
}

// Remove removes the code from the storage. All codes for the specified kind and principal are removed.
func (s *DefaultManager) Remove(ctx context.Context, kind Kind, principal string) bool {
	err := s.storage.Remove(ctx, kind, principal)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove OTP")
		return false
	}

	return true
}
