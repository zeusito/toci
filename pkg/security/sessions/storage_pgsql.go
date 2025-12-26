package sessions

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type PrincipalSessionModel struct {
	bun.BaseModel `bun:"table:user_sessions,alias:usersess"`
	ID            string          `bun:"id,pk"`
	PrincipalID   string          `bun:"principal_id"`
	IPAddress     string          `bun:"ip_address"`
	Metadata      PrincipalClaims `bun:"metadata"`
	ExpiresAt     time.Time       `bun:"expires_at"`
	CreatedAt     time.Time       `bun:"created_at"`
}

type PgSQLStorage struct {
	db *bun.DB
}

func NewPgSQLStorage(db *bun.DB) *PgSQLStorage {
	return &PgSQLStorage{db: db}
}

func (s *PgSQLStorage) Set(ctx context.Context, hashedToken string, claims PrincipalClaims, expiresAt time.Time) error {
	now := time.Now().UTC()

	data := &PrincipalSessionModel{
		ID:          hashedToken,
		PrincipalID: claims.PrincipalID,
		Metadata:    claims,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}

	_, err := s.db.NewInsert().Model(data).Exec(ctx)

	return err
}

func (s *PgSQLStorage) Get(ctx context.Context, hashedToken string) (PrincipalClaims, error) {
	var sessionData PrincipalSessionModel

	err := s.db.NewSelect().
		Model(&sessionData).
		Where("id = ?", hashedToken).
		Where("expires_at > ?", time.Now().UTC()).
		Scan(ctx, &sessionData)

	if err != nil {
		return PrincipalClaims{IsAuthenticated: false}, err
	}

	return sessionData.Metadata, nil
}

func (s *PgSQLStorage) Remove(ctx context.Context, hashedToken string) error {
	_, err := s.db.NewDelete().
		Model((*PrincipalSessionModel)(nil)).
		Where("id = ?", hashedToken).
		Exec(ctx)

	return err
}
