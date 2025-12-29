package security

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/toci/pkg/security/sessions"
)

// AuthenticationFilter is a middleware that checks if the request has a valid token
func AuthenticationFilter(sessionManager sessions.Manager) func(http.Handler) http.Handler {
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
			record, ok := sessionManager.GetSession(r.Context(), token)
			if !ok {
				log.Warn().Msg("session not authenticated")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims := sessions.ClaimsFromSession(record)

			// Add claims to context
			ctx := sessions.AddToContext(r.Context(), claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
