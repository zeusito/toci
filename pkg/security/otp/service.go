package otp

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type Kind string

const (
	KindPanelistPassword Kind = "panelist_password"
	KindEmployeePassword Kind = "employee_password"
)

// otpData internal struct used to store OTP data, not exposed to the outside world
type otpData struct {
	ID        string
	Kind      Kind
	Principal string
	ExpiresAt time.Time
}

type Manager interface {
	GenerateCode(ctx context.Context, length int, kind Kind, principal string) (string, bool)
	VerifyCode(ctx context.Context, kind Kind, principal string, code string) bool
	Remove(ctx context.Context, kind Kind, principal string) bool
}

type Storage interface {
	Put(ctx context.Context, kind Kind, principal, hashedCode string, expiresAt time.Time) error
	Get(ctx context.Context, kind Kind, principal string) (*otpData, error)
	Remove(ctx context.Context, kind Kind, principal string) error
}

func NewManagerWithPgSQLStorage(db *bun.DB, hasherSecret string) (Manager, bool) {
	theHasher, err := hasher.NewHmacSHA256(hasherSecret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create hasher")
		return nil, false
	}

	storage := NewPgSQLStore(db)

	return &DefaultManager{
		hashingAlgo:        theHasher,
		storage:            storage,
		expirationDuration: 5 * time.Minute,
	}, true
}
