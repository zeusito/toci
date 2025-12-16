package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeusito/toci/pkg/terrors"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5/middleware"
)

// RenderJSON is a helper function to write a JSON response
func RenderJSON(ctx context.Context, w http.ResponseWriter, httpStatusCode int, payload any) {
	// Headers
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(js)
}

// RenderError Renders an error with some sane defaults.
func RenderError(ctx context.Context, w http.ResponseWriter, err error) {
	var terrorToRender *terrors.Terror

	if !errors.As(err, &terrorToRender) {
		terrorToRender = terrors.Unknown(err.Error())
	}

	RenderJSON(ctx, w, terrorToRender.HttpStatusCode, terrorToRender)
}

func SimpleSuccessResponseBody() map[string]string {
	return map[string]string{
		"code":    "Success",
		"message": "action completed successfully",
	}
}
