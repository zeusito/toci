package handlers

import (
	"net/http"

	"github.com/zeusito/toci/pkg/router"

	"github.com/go-chi/chi/v5"
)

type HealthController struct{}

func NewHealthController(mux *chi.Mux) *HealthController {
	c := &HealthController{}

	mux.Get("/health/readiness", c.handleReadiness)
	mux.Get("/health/liveness", c.handleLiveness)

	return c
}

func (c *HealthController) handleReadiness(w http.ResponseWriter, req *http.Request) {
	router.RenderJSON(req.Context(), w, http.StatusOK, router.SimpleSuccessResponseBody())
}

func (c *HealthController) handleLiveness(w http.ResponseWriter, req *http.Request) {
	router.RenderJSON(req.Context(), w, http.StatusOK, router.SimpleSuccessResponseBody())
}
