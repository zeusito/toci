package otp

import (
	"github.com/zeusito/toci/pkg/toolbox"
	"github.com/zeusito/toci/pkg/toolbox/hasher"
)

type OneTimePassword struct {
	Secret       string
	HashedSecret string
}

type OneTimePasswordGenerator struct {
	hashingAlgo hasher.Hasher
}

func NewOneTimePasswordGenerator(hashingAlgo hasher.Hasher) *OneTimePasswordGenerator {
	return &OneTimePasswordGenerator{hashingAlgo: hashingAlgo}
}

func (g *OneTimePasswordGenerator) New(length int) (*OneTimePassword, error) {
	code := toolbox.SecureRandomString(length)

	hashedCode, err := g.hashingAlgo.Hash(code)
	if err != nil {
		return nil, err
	}

	return &OneTimePassword{
		Secret:       code,
		HashedSecret: hashedCode,
	}, nil
}

func (g *OneTimePasswordGenerator) Verify(hashedCode, code string) bool {
	return g.hashingAlgo.Verify(code, hashedCode)
}
