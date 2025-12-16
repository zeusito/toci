package sessions

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type OpaqueToken struct {
	tokenHasher hasher.Hasher
	prefix      string
	suffix      string
}

func (ot *OpaqueToken) String() string {
	return ot.prefix + "_" + ot.suffix
}

func (ot *OpaqueToken) SecureHashedString() (string, error) {
	return ot.tokenHasher.Hash(ot.String())
}

func (ot *OpaqueToken) Verify(hashed string) bool {
	return ot.tokenHasher.Verify(ot.String(), hashed)
}

// NewOpaqueToken generates a new opaque token following the spec:
// - Generate a UUID v4 (preferred due to its randomness versus v7)
// - Apply base64 URL-safe encoding
// - Prepend the prefix and an underscore
func NewOpaqueToken(prefix string, th hasher.Hasher) *OpaqueToken {
	baseId := uuid.New()
	encoded := base64.URLEncoding.EncodeToString(baseId[:])

	return &OpaqueToken{
		tokenHasher: th,
		prefix:      prefix,
		suffix:      encoded,
	}
}

// DecodeOpaqueTokenFromString decodes a string into an opaque token
func DecodeOpaqueTokenFromString(token string, th hasher.Hasher) (*OpaqueToken, error) {
	parts := strings.SplitN(token, "_", 2)

	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	return &OpaqueToken{
		tokenHasher: th,
		prefix:      parts[0],
		suffix:      parts[1],
	}, nil
}
