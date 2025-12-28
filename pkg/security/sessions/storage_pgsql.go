package sessions

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type PrincipalSessionRecord struct {
	bun.BaseModel `bun:"table:user_sessions,alias:usersess"`
	ID            string      `bun:"id,pk"`
	PrincipalID   string      `bun:"principal_id"`
	IPAddress     string      `bun:"ip_address"`
	Metadata      sessionData `bun:"metadata"`
	ExpiresAt     time.Time   `bun:"expires_at"`
	CreatedAt     time.Time   `bun:"created_at"`
}

type PgSQLStorage struct {
	db *bun.DB
}

func NewPgSQLStorage(db *bun.DB) *PgSQLStorage {
	return &PgSQLStorage{db: db}
}

func (s *PgSQLStorage) Set(ctx context.Context, hashedID string, data *sessionData) error {
	record := &PrincipalSessionRecord{
		ID:          hashedID,
		PrincipalID: data.PrincipalID,
		Metadata:    *data,
		ExpiresAt:   data.ExpiresAt,
		CreatedAt:   time.Now().UTC(),
	}

	_, err := s.db.NewInsert().Model(record).Exec(ctx)

	return err
}

func (s *PgSQLStorage) Get(ctx context.Context, hashedID string) (*sessionData, error) {
	var record PrincipalSessionRecord

	err := s.db.NewSelect().
		Model(&record).
		Where("id = ?", hashedID).
		Where("expires_at > ?", time.Now().UTC()).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &sessionData{
		PrincipalID: record.PrincipalID,
		OrgID:       record.Metadata.OrgID,
		Roles:       record.Metadata.Roles,
		ExpiresAt:   record.ExpiresAt,
		CreatedAt:   record.CreatedAt,
	}, nil
}

func (s *PgSQLStorage) Remove(ctx context.Context, hashedID string) error {
	_, err := s.db.NewDelete().
		Model((*PrincipalSessionRecord)(nil)).
		Where("id = ?", hashedID).
		Exec(ctx)

	return err
}
