package signin

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/internal/actions"
	"github.com/zeusito/toci/internal/dbmodels"
	"github.com/zeusito/toci/pkg/security/otp"
	"github.com/zeusito/toci/pkg/security/sessions"
	"github.com/zeusito/toci/pkg/terrors"
	"github.com/zeusito/toci/pkg/toolbox"
)

type DefaultService struct {
	repo           Repo
	otpManager     otp.Manager
	sessionManager sessions.Manager
	asyncActions   actions.Service
}

func NewDefaultService(repo Repo, otpManager otp.Manager, sessionManager sessions.Manager, asyncActions actions.Service) Service {
	return &DefaultService{repo: repo, otpManager: otpManager, sessionManager: sessionManager, asyncActions: asyncActions}
}

func (s *DefaultService) SignInWithEmailOTP(ctx context.Context, email string, source string) error {
	requestID := toolbox.GetRequestID(ctx)

	// Normalize email to lowercase
	email = strings.ToLower(email)

	log.Info().Str("trace", requestID).Msgf("login with email and password: %s", email)

	record, err := s.repo.FindOneByEmail(ctx, email)
	if err != nil {
		log.Warn().Str("trace", requestID).Msgf("failed to find user by email: %s", email)
		return terrors.UnAuthorized("credentials are invalid")
	}

	// Check if user is not active
	if record.Status != dbmodels.IdentityStatusActive {
		// Is it locked?
		if record.Status == dbmodels.IdentityStatusLocked {
			log.Warn().Str("trace", requestID).Msgf("user is locked: %s", email)
			return terrors.UnAuthorized("credentials are invalid")
		}

		log.Warn().Str("trace", requestID).Msgf("user is not active: %s", email)
		return terrors.UnAuthorized("credentials are invalid")
	}

	// Generate a one time password
	code, ok := s.otpManager.GenerateCode(ctx, 6, otp.CodeKindUserPassword, email)
	if !ok {
		log.Warn().Str("trace", requestID).Msgf("failed to generate one time password: %s", email)
		return terrors.UnAuthorized("credentials are invalid")
	}

	log.Info().Str("trace", requestID).Msg("one time password generated")

	s.asyncActions.SendOTPByEmail(ctx, code, email)

	return nil
}

func (s *DefaultService) VerifyEmailOTP(ctx context.Context, code, email string) (*SignInResponse, error) {
	requestID := toolbox.GetRequestID(ctx)
	now := time.Now().UTC()

	// Normalize email to lowercase
	email = strings.ToLower(email)

	log.Info().Str("trace", requestID).Msgf("verify email OTP: %s", email)

	// Verify the code
	ok := s.otpManager.VerifyCode(ctx, otp.CodeKindUserPassword, email, code)
	if !ok {
		log.Warn().Str("trace", requestID).Msgf("failed to verify code: %s", code)
		return nil, terrors.UnAuthorized("credentials are invalid")
	}

	// Retrieve the identity data
	record, err := s.repo.FindOneByEmail(ctx, email)
	if err != nil {
		log.Warn().Str("trace", requestID).Err(err).Msgf("failed to find identity: %s", email)
		return nil, terrors.UnAuthorized("credentials are invalid")
	}

	// Generate a session
	sessionData := sessions.Session{
		PrincipalID: record.ID,
		Metadata: sessions.SessionMetadata{
			"roles": "user",
		},
		ExpiresAt: now.Add(time.Hour * 24),
		CreatedAt: now,
	}
	sessionID, ok := s.sessionManager.CreateSession(ctx, sessionData, sessionData.ExpiresAt)
	if !ok {
		log.Warn().Str("trace", requestID).Msgf("failed to create session: %s", email)
		return nil, terrors.UnAuthorized("credentials are invalid")
	}

	// Return the session ID
	return &SignInResponse{
		AccessToken: sessionID,
		TokenType:   "Bearer",
		ExpiresAt:   sessionData.ExpiresAt,
	}, nil

}

func (s *DefaultService) SignInWithOpenID(ctx context.Context, provider, token string, source string) (*SignInResponse, error) {
	return nil, nil
}
