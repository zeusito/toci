package dbmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type IdentityStatus string

const (
	IdentityStatusActive    IdentityStatus = "active"
	IdentityStatusDeleted   IdentityStatus = "deleted"
	IdentityStatusSuspended IdentityStatus = "suspended"
	IdentityStatusLocked    IdentityStatus = "locked"
)

type IdentityRecord struct {
	bun.BaseModel       `bun:"table:identities,alias:i"`
	ID                  string         `bun:"id,pk"`
	Email               string         `bun:"email"`
	FirstName           string         `bun:"first_name"`
	LastName            string         `bun:"last_name"`
	Status              IdentityStatus `bun:"status"`
	EmailVerifiedAt     *time.Time     `bun:"email_verified_at"`
	FailedLoginAttempts int            `bun:"failed_login_attempts"`
	LockExpiresAt       time.Time      `bun:"lock_expires_at"`
	LastLoginAt         time.Time      `bun:"last_login_at"`
	CreatedAt           time.Time      `bun:"created_at"`
	UpdatedAt           time.Time      `bun:"updated_at"`
}
