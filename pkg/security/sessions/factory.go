package sessions

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

// sessionData internal use only, not exposed to the outside world. Customize as needed
type sessionData struct {
	PrincipalID string
	OrgID       string
	Roles       []string
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type Manager interface {
	NewSession(ctx context.Context, data PrincipalClaims, validForDuration time.Duration) (string, bool)
	GetSession(ctx context.Context, token string) PrincipalClaims
	RemoveSession(ctx context.Context, token string) bool
	CleanUpExpiredSessions(ctx context.Context)
}

type Storage interface {
	Set(ctx context.Context, hashedID string, data *sessionData) error
	Get(ctx context.Context, hashedID string) (*sessionData, error)
	Remove(ctx context.Context, hashedID string) error
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
