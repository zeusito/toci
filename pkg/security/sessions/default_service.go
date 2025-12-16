package sessions

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type DefaultService struct {
	storage     Storage
	tokenHasher hasher.Hasher
	tokenPrefix string
}

func NewDefaultService(store Storage, th hasher.Hasher, tokenPrefix string) *DefaultService {
	return &DefaultService{
		storage:     store,
		tokenHasher: th,
		tokenPrefix: tokenPrefix,
	}
}

func (s *DefaultService) NewSession(ctx context.Context, data *AuthClaims, validForDuration time.Duration) (string, error) {
	log.Info().Msgf("Creating new session for principal %s", data.Principal)

	// Create a new opaque token
	opaqueToken := NewOpaqueToken(s.tokenPrefix, s.tokenHasher)
	hashedToken, err := opaqueToken.SecureHashedString()

	if err != nil {
		return "", err
	}

	// Persist the session in storage
	expiresAt := time.Now().UTC().Add(validForDuration)

	err = s.storage.Set(ctx, hashedToken, data, expiresAt)
	if err != nil {
		return "", err
	}

	// Return the opaque token string
	return opaqueToken.String(), nil
}

func (s *DefaultService) GetSession(ctx context.Context, token string) *AuthClaims {
	log.Info().Msgf("Getting session from token...")

	opaqueToken, _ := DecodeOpaqueTokenFromString(token, s.tokenHasher)
	hashedToken, err := opaqueToken.SecureHashedString()

	if err != nil {
		log.Warn().Err(err).Msg("Failed to decode token")
		return &AuthClaims{IsAuthenticated: false}
	}

	session, err := s.storage.Get(ctx, hashedToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get session from storage")
		return &AuthClaims{IsAuthenticated: false}
	}

	return session
}

func (s *DefaultService) RemoveSession(ctx context.Context, token string) error {
	log.Info().Msgf("Removing session...")

	// Decode the opaque token
	opaqueToken, err := DecodeOpaqueTokenFromString(token, s.tokenHasher)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to decode token")
		return errors.New("invalid token")
	}

	// Extract the hashed token
	hashedToken, err := opaqueToken.SecureHashedString()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to decode token")
		return errors.New("invalid token")
	}

	// Remove the session from storage
	err = s.storage.Remove(ctx, hashedToken)

	return err
}
