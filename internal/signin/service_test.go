package signin

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusito/toci/internal/actions"
	"github.com/zeusito/toci/internal/dbmodels"
	"github.com/zeusito/toci/pkg/security/otp"
	"github.com/zeusito/toci/pkg/security/sessions"
)

func TestSignInWithEmailOTPInvalidEmail(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepo(t)
	ottManager := otp.NewMockManager(t)
	sessionManager := sessions.NewMockManager(t)
	asyncActions := actions.NewMockService(t)

	svc := NewDefaultService(repo, ottManager, sessionManager, asyncActions)

	// Expectations
	repo.EXPECT().FindOneByEmail(ctx, "none@my.com").Return(nil, errors.New("record not found"))

	err := svc.SignInWithEmailOTP(ctx, "none@my.com", "web")
	assert.Error(t, err, "expected error for non-existent email")
}

func TestSignInWithEmailOTPUserLocked(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepo(t)
	ottManager := otp.NewMockManager(t)
	sessionManager := sessions.NewMockManager(t)
	asyncActions := actions.NewMockService(t)

	svc := NewDefaultService(repo, ottManager, sessionManager, asyncActions)

	// Expectations
	repo.EXPECT().FindOneByEmail(ctx, "none@my.com").Return(&dbmodels.IdentityRecord{
		ID:     "1",
		Email:  "none@my.com",
		Status: dbmodels.IdentityStatusLocked,
	}, nil)

	err := svc.SignInWithEmailOTP(ctx, "none@my.com", "web")
	assert.Error(t, err, "expected error for locked user")
}

func TestSignInWithEmailOTPFailToGenerateOTP(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepo(t)
	ottManager := otp.NewMockManager(t)
	sessionManager := sessions.NewMockManager(t)
	asyncActions := actions.NewMockService(t)

	svc := NewDefaultService(repo, ottManager, sessionManager, asyncActions)

	// Expectations
	repo.EXPECT().FindOneByEmail(ctx, "none@my.com").Return(&dbmodels.IdentityRecord{
		ID:     "1",
		Email:  "none@my.com",
		Status: dbmodels.IdentityStatusActive,
	}, nil)

	ottManager.EXPECT().GenerateCode(ctx, 6, otp.CodeKindUserPassword, "none@my.com").Return("", false)

	err := svc.SignInWithEmailOTP(ctx, "none@my.com", "web")
	assert.Error(t, err, "expected error for failed to generate OTP")
}

func TestSignInWithEmailOTPSuccess(t *testing.T) {
	ctx := context.Background()
	repo := NewMockRepo(t)
	ottManager := otp.NewMockManager(t)
	sessionManager := sessions.NewMockManager(t)
	asyncActions := actions.NewMockService(t)

	svc := NewDefaultService(repo, ottManager, sessionManager, asyncActions)

	// Expectations
	repo.EXPECT().FindOneByEmail(ctx, "none@my.com").Return(&dbmodels.IdentityRecord{
		ID:     "1",
		Email:  "none@my.com",
		Status: dbmodels.IdentityStatusActive,
	}, nil)

	ottManager.EXPECT().GenerateCode(ctx, 6, otp.CodeKindUserPassword, "none@my.com").Return("123456", true)

	asyncActions.EXPECT().SendOTPByEmail(ctx, "123456", "none@my.com")

	err := svc.SignInWithEmailOTP(ctx, "none@my.com", "web")
	assert.NoError(t, err, "expected no error for successful OTP generation")
}
