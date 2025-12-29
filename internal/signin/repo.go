package signin

import (
	"context"

	"github.com/zeusito/toci/internal/dbmodels"
)

type Repo interface {
	FindOneByEmail(ctx context.Context, email string) (*dbmodels.IdentityRecord, error)
}
