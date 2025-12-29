package signin

import (
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"github.com/zeusito/toci/internal/actions"
	"github.com/zeusito/toci/pkg/security/otp"
	"github.com/zeusito/toci/pkg/security/sessions"
)

func InitModule(mux *chi.Mux, db *bun.DB, optManager otp.Manager, sessionManager sessions.Manager, asyncActions actions.Service) {
	repo := NewDefaultRepo(db)
	svc := NewDefaultService(repo, optManager, sessionManager, asyncActions)
	_ = NewController(mux, svc)
}
