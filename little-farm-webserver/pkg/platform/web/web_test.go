package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"little-farm-webserver/pkg/apperr"

	"github.com/stretchr/testify/assert"
)

func decodeGzipJSON(t *testing.T, body []byte, target any) {
	t.Helper()
	gr, err := gzip.NewReader(bytes.NewReader(body))
	assert.NoError(t, err)
	defer gr.Close()
	assert.NoError(t, json.NewDecoder(gr).Decode(target))
}

func TestRespondWithJSON(t *testing.T) {
	t.Run("It should write the given status code", func(t *testing.T) {
		rec := httptest.NewRecorder()
		RespondWithJSON(rec, http.StatusCreated, map[string]string{"key": "value"})
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("It should set Content-Type to application/json", func(t *testing.T) {
		rec := httptest.NewRecorder()
		RespondWithJSON(rec, http.StatusOK, map[string]string{})
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("It should set Content-Encoding to gzip", func(t *testing.T) {
		rec := httptest.NewRecorder()
		RespondWithJSON(rec, http.StatusOK, map[string]string{})
		assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))
	})

	t.Run("It should set Content-Length to the byte size of the gzip body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		RespondWithJSON(rec, http.StatusOK, map[string]string{"key": "value"})
		assert.NotEmpty(t, rec.Header().Get("Content-Length"))
	})

	t.Run("It should encode the data as gzip-compressed JSON in the body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		RespondWithJSON(rec, http.StatusOK, map[string]string{"hello": "world"})
		var result map[string]string
		decodeGzipJSON(t, rec.Body.Bytes(), &result)
		assert.Equal(t, map[string]string{"hello": "world"}, result)
	})
}

func TestRespondWithError(t *testing.T) {
	t.Run("It should respond with 404 when error wraps ErrNotFound", func(t *testing.T) {
		rec := httptest.NewRecorder()
		err := fmt.Errorf("%w: thing not found", apperr.ErrNotFound)
		RespondWithError(rec, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		var body map[string]string
		decodeGzipJSON(t, rec.Body.Bytes(), &body)
		assert.Equal(t, err.Error(), body["error"])
	})

	t.Run("It should respond with 409 when error wraps ErrConflict", func(t *testing.T) {
		rec := httptest.NewRecorder()
		err := fmt.Errorf("%w: duplicate entry", apperr.ErrConflict)
		RespondWithError(rec, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		var body map[string]string
		decodeGzipJSON(t, rec.Body.Bytes(), &body)
		assert.Equal(t, err.Error(), body["error"])
	})

	t.Run("It should respond with 400 when error wraps ErrValidation", func(t *testing.T) {
		rec := httptest.NewRecorder()
		err := fmt.Errorf("%w: bad input", apperr.ErrValidation)
		RespondWithError(rec, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var body map[string]string
		decodeGzipJSON(t, rec.Body.Bytes(), &body)
		assert.Equal(t, err.Error(), body["error"])
	})

	t.Run("It should respond with 422 when error wraps ErrUnprocessableEntity", func(t *testing.T) {
		rec := httptest.NewRecorder()
		err := fmt.Errorf("%w: cannot process", apperr.ErrUnprocessableEntity)
		RespondWithError(rec, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		var body map[string]string
		decodeGzipJSON(t, rec.Body.Bytes(), &body)
		assert.Equal(t, err.Error(), body["error"])
	})

	t.Run("It should respond with 500 and a generic message for unknown errors", func(t *testing.T) {
		rec := httptest.NewRecorder()
		RespondWithError(rec, fmt.Errorf("some unexpected error"))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var body map[string]string
		decodeGzipJSON(t, rec.Body.Bytes(), &body)
		assert.Equal(t, apperr.ErrInternal.Error(), body["error"])
	})
}
