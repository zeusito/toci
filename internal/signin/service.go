package signin

import "context"

type Service interface {
	SignInWithEmailOTP(ctx context.Context, email string, source string) error
	VerifyEmailOTP(ctx context.Context, code, email string) (*SignInResponse, error)
	SignInWithOpenID(ctx context.Context, provider, token string, source string) (*SignInResponse, error)
}
