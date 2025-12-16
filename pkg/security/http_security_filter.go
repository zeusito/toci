package security

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/pkg/security/sessions"
)

// AuthenticationFilter is a middleware that checks if the request has a valid token
func AuthenticationFilter(sessionSvc sessions.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				log.Warn().Msg("no token provided")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Remove the "Bearer " prefix
			token = token[7:]

			// Validate the token
			claims := sessionSvc.GetSession(r.Context(), token)
			if !claims.IsAuthenticated {
				log.Warn().Msg("invalid token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := sessions.AddAuthClaimsToContext(r.Context(), claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
