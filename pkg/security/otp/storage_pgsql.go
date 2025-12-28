package otp

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type OneTimeTokenRecord struct {
	bun.BaseModel `bun:"table:user_otts,alias:ott"`
	ID            string    `bun:"id,pk"`
	Kind          Kind      `bun:"kind"`
	Principal     string    `bun:"principal"`
	ExpiresAt     time.Time `bun:"expires_at"`
	CreatedAt     time.Time `bun:"created_at"`
}

type PgSQLStore struct {
	db *bun.DB
}

func NewPgSQLStore(db *bun.DB) *PgSQLStore {
	return &PgSQLStore{
		db: db,
	}
}

// Put stores a new OTP in the database
func (s *PgSQLStore) Put(ctx context.Context, kind Kind, principal, hashedCode string, expiresAt time.Time) error {
	_, err := s.db.NewInsert().
		Model(&OneTimeTokenRecord{
			ID:        hashedCode,
			Kind:      kind,
			Principal: principal,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now().UTC(),
		}).
		Exec(ctx)

	return err
}

// Get retrieves an OTP from the database, by default only the latest one for the given kind and principal is returned
func (s *PgSQLStore) Get(ctx context.Context, kind Kind, principal string) (*otpData, error) {
	var model OneTimeTokenRecord

	err := s.db.NewSelect().
		Model(&model).
		Where("principal = ?", principal).
		Where("kind = ?", kind).
		Where("expires_at > ?", time.Now().UTC()).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &otpData{
		ID:        model.ID,
		Kind:      model.Kind,
		Principal: model.Principal,
		ExpiresAt: model.ExpiresAt,
	}, nil
}

// Remove deletes all codes for a given kind and principal
func (s *PgSQLStore) Remove(ctx context.Context, kind Kind, principal string) error {
	_, err := s.db.NewDelete().
		Model((*OneTimeTokenRecord)(nil)).
		Where("principal = ?", principal).
		Where("kind = ?", kind).
		Exec(ctx)

	return err
}
