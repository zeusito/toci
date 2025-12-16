package sessions

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type IdentitySessionModel struct {
	bun.BaseModel  `bun:"table:identity_sessions,alias:idsess"`
	ID             string    `bun:"id,pk"`
	IdentityID     string    `bun:"identity_id"`
	OrganizationID string    `bun:"organization_id"`
	Roles          []string  `bun:"roles"`
	IpAddress      string    `bun:"ip_address"`
	ExpiresAt      time.Time `bun:"expires_at"`
	CreatedAt      time.Time `bun:"created_at"`
}

type DBStorage struct {
	db *bun.DB
}

func NewDBStore(db *bun.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (s *DBStorage) Set(ctx context.Context, token string, claims *AuthClaims, expiresAt time.Time) error {
	now := time.Now().UTC()
	data := &IdentitySessionModel{
		ID:             token,
		IdentityID:     claims.Principal,
		OrganizationID: claims.Metadata.OrganizationID,
		Roles:          claims.Roles,
		ExpiresAt:      expiresAt,
		CreatedAt:      now,
	}

	_, err := s.db.NewInsert().Model(data).Exec(ctx)

	return err
}

func (s *DBStorage) Get(ctx context.Context, token string) (*AuthClaims, error) {
	var sessionData IdentitySessionModel

	err := s.db.NewSelect().
		Model(&sessionData).
		Where("id = ?", token).
		Where("expires_at > ?", time.Now().UTC()).
		Scan(ctx, &sessionData)

	if err != nil {
		return nil, err
	}

	claims := &AuthClaims{
		IsAuthenticated: true,
		Principal:       sessionData.IdentityID,
		Roles:           sessionData.Roles,
		Metadata: AuthClaimsMetadata{
			OrganizationID: sessionData.OrganizationID,
		},
	}

	return claims, nil
}

func (s *DBStorage) Remove(ctx context.Context, token string) error {
	_, err := s.db.NewDelete().
		Model((*IdentitySessionModel)(nil)).
		Where("id = ?", token).
		Exec(ctx)

	return err
}
