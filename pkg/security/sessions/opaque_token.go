package sessions

import (
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

// NewOpaqueToken generates a new opaque token following the spec:
// - Generate a UUID v4 (preferred due to its randomness versus v7)
// - Apply base64 URL-safe encoding
func NewOpaqueToken(th hasher.Hasher) (token string, hashedToken string, err error) {
	baseId := uuid.New()
	encoded := base64.URLEncoding.EncodeToString(baseId[:])

	ht, err := th.Hash(encoded)

	return encoded, ht, err
}
