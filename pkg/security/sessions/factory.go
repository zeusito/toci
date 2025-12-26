package sessions

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type Manager interface {
	NewSession(ctx context.Context, data PrincipalClaims, validForDuration time.Duration) (string, bool)
	GetSession(ctx context.Context, token string) PrincipalClaims
	RemoveSession(ctx context.Context, token string) bool
	CleanUpExpiredSessions(ctx context.Context)
}

type Storage interface {
	Set(ctx context.Context, hashedToken string, claims PrincipalClaims, expiresAt time.Time) error
	Get(ctx context.Context, hashedToken string) (PrincipalClaims, error)
	Remove(ctx context.Context, hashedToken string) error
}

func NewManagerWithPgSQL(db *bun.DB, hasherSecret string) (Manager, bool) {
	theHasher, err := hasher.NewHmacSHA256(hasherSecret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create hasher")
		return nil, false
	}

	return &DefaultManager{
		storage:     NewPgSQLStorage(db),
		tokenHasher: theHasher,
	}, true
}
