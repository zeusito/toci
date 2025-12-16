package toolbox

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/go-chi/chi/v5/middleware"
)

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	return middleware.GetReqID(ctx)
}

// SecureRandomString generates a cryptographically secure random string of the specified length.
// It uses characters from the set [a-zA-Z0-9].
func SecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	maxNumber := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, maxNumber)
		if err != nil {
			return ""
		}
		b[i] = charset[n.Int64()]
	}

	return string(b)
}
