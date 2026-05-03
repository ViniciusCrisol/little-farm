package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"little-farm-webserver/pkg/apperr"
)

func RespondWithError(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		RespondWithJSON(response, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, apperr.ErrConflict):
		RespondWithJSON(response, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, apperr.ErrValidation):
		RespondWithJSON(response, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, apperr.ErrUnprocessableEntity):
		RespondWithJSON(response, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		RespondWithJSON(response, http.StatusInternalServerError, map[string]string{"error": apperr.ErrInternal.Error()})
	}
}

func RespondWithJSON(response http.ResponseWriter, status int, data any) {
	var buff bytes.Buffer
	gz := gzip.NewWriter(&buff)
	if err := json.NewEncoder(gz).Encode(data); err != nil {
		slog.Error("failed to encode response", slog.Any("data", data), slog.String("error", err.Error()))
		return
	}
	gz.Close()

	contentLen := strconv.Itoa(buff.Len())
	response.Header().Set("Content-Encoding", "gzip")
	response.Header().Set("Content-Length", contentLen)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	if _, err := response.Write(buff.Bytes()); err != nil {
		slog.Error("failed to write response", slog.Any("data", data), slog.String("error", err.Error()))
	}
}
