package signin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zeusito/toci/pkg/router"
)

type Controller struct {
	svc Service
}

func NewController(mux *chi.Mux, svc Service) *Controller {
	c := &Controller{svc: svc}

	mux.Post("/v1/auth/otp/login", c.handleLogin)
	mux.Post("/v1/auth/otp/verify", c.handleVerifyOTP)
	mux.Post("/v1/auth/oidc/callback", c.handleOIDCLogin)

	return c
}

func (c *Controller) handleLogin(w http.ResponseWriter, req *http.Request) {
	var body LoginWithEmailOTPRequest
	err := router.BindBody(req, &body)
	if err != nil {
		router.RenderError(req.Context(), w, err)
		return
	}

	err = c.svc.SignInWithEmailOTP(req.Context(), body.Email, body.Source)

	if err != nil {
		router.RenderError(req.Context(), w, err)
		return
	}

	router.RenderJSON(req.Context(), w, http.StatusOK, router.SimpleSuccessResponseBody())
}

func (c *Controller) handleVerifyOTP(w http.ResponseWriter, req *http.Request) {
	var body VerifyEmailOTPRequest
	err := router.BindBody(req, &body)
	if err != nil {
		router.RenderError(req.Context(), w, err)
		return
	}

	resp, err := c.svc.VerifyEmailOTP(req.Context(), body.Code, body.Email)
	if err != nil {
		router.RenderError(req.Context(), w, err)
		return
	}

	router.RenderJSON(req.Context(), w, http.StatusOK, resp)
}

func (c *Controller) handleOIDCLogin(w http.ResponseWriter, req *http.Request) {
	var body OIDCLoginRequest
	err := router.BindBody(req, &body)
	if err != nil {
		router.RenderError(req.Context(), w, err)
		return
	}

	resp, err := c.svc.SignInWithOpenID(req.Context(), body.Provider, body.Token, body.Source)

	if err != nil {
		router.RenderError(req.Context(), w, err)
		return
	}

	router.RenderJSON(req.Context(), w, http.StatusOK, resp)
}
