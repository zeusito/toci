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

func (s *DefaultManager) CreateSession(ctx context.Context, data Session, expiresAt time.Time) (string, bool) {
	log.Info().Msgf("Creating new session for principal %s", data.PrincipalID)

	// Create a new opaque token
	token, hashedToken, err := NewOpaqueToken(s.tokenHasher)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create new opaque token")
		return "", false
	}

	// Persist the session in storage
	sessionData := &Session{
		PrincipalID: data.PrincipalID,
		Metadata:    data.Metadata,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now().UTC(),
	}

	err = s.storage.Set(ctx, hashedToken, sessionData)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to persist session in storage")
		return "", false
	}

	// Return the opaque token string
	return token, true
}

func (s *DefaultManager) GetSession(ctx context.Context, token string) (*Session, bool) {
	log.Info().Msgf("Getting session from token...")

	hashedToken, err := s.tokenHasher.Hash(token)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to decode token")
		return nil, false
	}

	record, err := s.storage.Get(ctx, hashedToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get session from storage")
		return nil, false
	}

	// verify if the session is expired
	if record.ExpiresAt.Before(time.Now().UTC()) {
		log.Warn().Msg("Session is expired")
		return nil, false
	}

	return record, true
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
