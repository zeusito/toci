package sessions

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type SessionMetadata map[string]string

type Session struct {
	PrincipalID string
	Metadata    SessionMetadata
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type Manager interface {
	CreateSession(ctx context.Context, data Session, expiresAt time.Time) (string, bool)
	GetSession(ctx context.Context, token string) (*Session, bool)
	RemoveSession(ctx context.Context, token string) bool
	CleanUpExpiredSessions(ctx context.Context)
}

type Storage interface {
	Set(ctx context.Context, hashedID string, data *Session) error
	Get(ctx context.Context, hashedID string) (*Session, error)
	Remove(ctx context.Context, hashedID string) error
}

func NewManagerWithPgSQLStorage(db *bun.DB, hasherSecret string) (Manager, bool) {
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
