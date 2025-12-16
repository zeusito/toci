package sessions

import (
	"context"
	"time"
)

type Service interface {
	NewSession(ctx context.Context, data *AuthClaims, validForDuration time.Duration) (string, error)
	GetSession(ctx context.Context, token string) *AuthClaims
	RemoveSession(ctx context.Context, token string) error
}

type Storage interface {
	Set(ctx context.Context, token string, claims *AuthClaims, expiresAt time.Time) error
	Get(ctx context.Context, token string) (*AuthClaims, error)
	Remove(ctx context.Context, token string) error
}
