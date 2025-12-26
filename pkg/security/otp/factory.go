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

// OneTimePassword is a struct that represents an OTP, customize it as needed
type OneTimePassword struct {
	Kind       Kind
	Code       string
	HashedCode string
	ExpiresAt  time.Time
}

type Manager interface {
	GenerateOTP(ctx context.Context, length int, kind Kind) (OneTimePassword, bool)
	Retrieve(ctx context.Context, kind Kind, code string) (OneTimePassword, bool)
	Remove(ctx context.Context, code string) bool
}

type Storage interface {
	Put(ctx context.Context, data OneTimePassword, expiresAt time.Time) error
	Get(ctx context.Context, kind Kind, hashedCode string) (OneTimePassword, error)
	Remove(ctx context.Context, hashedCode string) error
}

func NewManagerWithPgSQL(db *bun.DB, hasherSecret string) (Manager, bool) {
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
