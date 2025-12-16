package terrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreconditionFailed(t *testing.T) {
	err := PreconditionFailed("test")
	assert.Equal(t, "test", err.ErrMessage, "PreconditionFailed should return the correct message")
	assert.Equal(t, "PreconditionFailed", err.ErrCode, "PreconditionFailed should return the correct code")
	assert.Equal(t, http.StatusBadRequest, err.HttpStatusCode, "PreconditionFailed should return the correct http status code")
}

func TestForbidden(t *testing.T) {
	err := Forbidden("test")
	assert.Equal(t, "test", err.ErrMessage, "Forbidden should return the correct message")
	assert.Equal(t, "ActionForbidden", err.ErrCode, "Forbidden should return the correct code")
	assert.Equal(t, http.StatusForbidden, err.HttpStatusCode, "Forbidden should return the correct http status code")
}

func TestRecordNotFound(t *testing.T) {
	err := RecordNotFound("test")
	assert.Equal(t, "test", err.ErrMessage, "RecordNotFound should return the correct message")
	assert.Equal(t, "RecordNotFound", err.ErrCode, "RecordNotFound should return the correct code")
	assert.Equal(t, http.StatusBadRequest, err.HttpStatusCode, "RecordNotFound should return the correct http status code")
}

func TestUnknown(t *testing.T) {
	err := Unknown("test")
	assert.Equal(t, "test", err.ErrMessage, "Unknown should return the correct message")
	assert.Equal(t, "UnknownError", err.ErrCode, "Unknown should return the correct code")
	assert.Equal(t, http.StatusInternalServerError, err.HttpStatusCode, "Unknown should return the correct http status code")
}

func TestTypeAssertion(t *testing.T) {
	var err error = PreconditionFailed("test")

	var terr *Terror
	ok := errors.As(err, &terr)

	assert.NotNil(t, terr)
	assert.True(t, ok)
	assert.Equal(t, "PreconditionFailed", terr.ErrCode)
	assert.Equal(t, "test", terr.ErrMessage)
}
