package sessions

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type DefaultManager struct {
	storage     Storage
	tokenHasher hasher.Hasher
}

func (s *DefaultManager) NewSession(ctx context.Context, data PrincipalClaims, validForDuration time.Duration) (string, bool) {
	log.Info().Msgf("Creating new session for principal %s", data.PrincipalID)

	// Create a new opaque token
	token, hashedToken, err := NewOpaqueToken(s.tokenHasher)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create new opaque token")
		return "", false
	}

	// Persist the session in storage
	expiresAt := time.Now().UTC().Add(validForDuration)

	err = s.storage.Set(ctx, hashedToken, data, expiresAt)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to persist session in storage")
		return "", false
	}

	// Return the opaque token string
	return token, true
}

func (s *DefaultManager) GetSession(ctx context.Context, token string) PrincipalClaims {
	log.Info().Msgf("Getting session from token...")

	hashedToken, err := s.tokenHasher.Hash(token)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to decode token")
		return PrincipalClaims{IsAuthenticated: false}
	}

	session, err := s.storage.Get(ctx, hashedToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get session from storage")
		return PrincipalClaims{IsAuthenticated: false}
	}

	return session
}

func (s *DefaultManager) RemoveSession(ctx context.Context, token string) bool {
	log.Info().Msgf("Removing session...")

	hashedToken, err := s.tokenHasher.Hash(token)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to decode token")
		return false
	}

	// Remove the session from storage
	err = s.storage.Remove(ctx, hashedToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to remove session from storage")
		return false
	}

	return true
}

func (s *DefaultManager) CleanUpExpiredSessions(ctx context.Context) {
	log.Info().Msg("Cleaning up expired sessions...")

	// TODO: Implement clean up expired sessions
}
