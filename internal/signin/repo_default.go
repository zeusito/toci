package signin

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/zeusito/toci/internal/dbmodels"
)

type defaultRepo struct {
	db *bun.DB
}

func NewDefaultRepo(db *bun.DB) Repo {
	return &defaultRepo{db: db}
}

func (r *defaultRepo) FindOneByEmail(ctx context.Context, email string) (*dbmodels.IdentityRecord, error) {
	return nil, nil
}
