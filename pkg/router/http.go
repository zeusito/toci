package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zeusito/toci/pkg/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type HTTPRouter struct {
	Mux *chi.Mux
	srv *http.Server
}

func NewHTTPRouter(cfgs config.ServerConfigurations) *HTTPRouter {
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request models (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(20 * time.Second))

	// Customizing the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfgs.Port),
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	return &HTTPRouter{
		Mux: router,
		srv: srv,
	}
}

func (s *HTTPRouter) Start() {
	log.Info().Msgf("Server listening on port %s", s.srv.Addr)
	_ = s.srv.ListenAndServe()
}

func (s *HTTPRouter) Shutdown(ctx context.Context) {
	log.Info().Msg("Server shutting down...")
	_ = s.srv.Shutdown(ctx)
}
