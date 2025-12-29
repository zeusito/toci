package sessions

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// PrincipalSessionRecord the database model for a session
type PrincipalSessionRecord struct {
	bun.BaseModel `bun:"table:user_sessions,alias:us"`
	ID            string          `bun:"id,pk"` // hashed ID
	PrincipalID   string          `bun:"principal_id"`
	IPAddress     string          `bun:"ip_address"`
	Metadata      SessionMetadata `bun:"metadata"`
	ExpiresAt     time.Time       `bun:"expires_at"`
	CreatedAt     time.Time       `bun:"created_at"`
}

type PgSQLStorage struct {
	db *bun.DB
}

func NewPgSQLStorage(db *bun.DB) Storage {
	return &PgSQLStorage{db: db}
}

// Set stores a session in the database
func (s *PgSQLStorage) Set(ctx context.Context, hashedID string, data *Session) error {
	record := &PrincipalSessionRecord{
		ID:          hashedID,
		PrincipalID: data.PrincipalID,
		Metadata:    data.Metadata,
		ExpiresAt:   data.ExpiresAt,
		CreatedAt:   data.CreatedAt,
	}

	_, err := s.db.NewInsert().Model(record).Exec(ctx)

	return err
}

// Get retrieves a session from the database, if it exists and is not expired
func (s *PgSQLStorage) Get(ctx context.Context, hashedID string) (*Session, error) {
	var record PrincipalSessionRecord

	err := s.db.NewSelect().
		Model(&record).
		Where("id = ?", hashedID).
		Where("expires_at > ?", time.Now().UTC()).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &Session{
		PrincipalID: record.PrincipalID,
		Metadata:    record.Metadata,
		ExpiresAt:   record.ExpiresAt,
		CreatedAt:   record.CreatedAt,
	}, nil
}

// Remove removes a session from the database
func (s *PgSQLStorage) Remove(ctx context.Context, hashedID string) error {
	_, err := s.db.NewDelete().
		Model((*PrincipalSessionRecord)(nil)).
		Where("id = ?", hashedID).
		Exec(ctx)

	return err
}
