package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusito/toci/pkg/terrors"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5/middleware"
)

func TestRenderJSON_Success(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	payload := map[string]string{"message": "hello"}
	expectedStatusCode := http.StatusOK

	RenderJSON(ctx, recorder, expectedStatusCode, payload)

	assert.Equal(t, expectedStatusCode, recorder.Code, "Status code should match expected value")
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Content type should be application/json")
	assert.Equal(t, "test-request-id", recorder.Header().Get(middleware.RequestIDHeader), "Request ID should be set")

	var responseBody map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)

	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Equal(t, payload["message"], responseBody["message"], "Response body should match expected payload")
}

func TestRenderJSON_MarshalError(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.Background()
	// Functions cannot be marshalled to JSON, this will cause an error
	payload := func() {}
	expectedStatusCode := http.StatusInternalServerError

	RenderJSON(ctx, recorder, http.StatusOK, payload) // StatusOK is intentional to test override

	assert.Equal(t, expectedStatusCode, recorder.Code, "Status code should match expected value")
	assert.True(t, len(recorder.Body.String()) > 0, "Response body should contain an error message")
}

func TestRenderError_WithTerror(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	terror := terrors.RecordNotFound("resource not found")
	expectedStatusCode := http.StatusBadRequest

	RenderError(ctx, recorder, terror)

	assert.Equal(t, expectedStatusCode, recorder.Code, "Status code should match expected value")
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Content type should be application/json")
	assert.Equal(t, "test-request-id", recorder.Header().Get(middleware.RequestIDHeader), "Request ID should be set")
}

func TestRenderError_WithGenericError(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "another-id")
	genericErr := errors.New("a generic error occurred")
	expectedStatusCode := http.StatusInternalServerError // Default for unknown errors

	RenderError(ctx, recorder, genericErr)

	assert.Equal(t, expectedStatusCode, recorder.Code, "Status code should match expected value")
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Content type should be application/json")
	assert.Equal(t, "another-id", recorder.Header().Get(middleware.RequestIDHeader), "Request ID should be set")
}

func TestSimpleSuccessResponseBody(t *testing.T) {
	expected := map[string]string{
		"code":    "Success",
		"message": "action completed successfully",
	}
	result := SimpleSuccessResponseBody()

	assert.Equal(t, expected, result, "Expected response body to match expected value")
}
