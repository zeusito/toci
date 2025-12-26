package otp

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type OneTimeTokenModel struct {
	bun.BaseModel `bun:"table:otts,alias:otts"`
	ID            string    `bun:"id,pk"`
	Kind          Kind      `bun:"kind"`
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

func (s *PgSQLStore) Put(ctx context.Context, data OneTimePassword, expiresAt time.Time) error {
	_, err := s.db.NewInsert().
		Model(&OneTimeTokenModel{
			ID:        data.HashedCode,
			Kind:      data.Kind,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now().UTC(),
		}).
		Exec(ctx)

	return err
}

func (s *PgSQLStore) Get(ctx context.Context, kind Kind, hashedCode string) (OneTimePassword, error) {
	var model OneTimeTokenModel

	err := s.db.NewSelect().
		Model(&model).
		Where("id = ?", hashedCode).
		Where("kind = ?", kind).
		Where("expires_at > ?", time.Now().UTC()).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return OneTimePassword{}, err
	}

	return OneTimePassword{
		Kind:       model.Kind,
		Code:       "", // we don't know this value here
		HashedCode: model.ID,
	}, nil
}

func (s *PgSQLStore) Remove(ctx context.Context, hashedCode string) error {
	_, err := s.db.NewDelete().
		Model((*OneTimeTokenModel)(nil)).
		Where("id = ?", hashedCode).
		Exec(ctx)

	return err
}
