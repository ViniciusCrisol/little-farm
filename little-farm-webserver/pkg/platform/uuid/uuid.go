package uuid

import (
	"log/slog"

	"github.com/google/uuid"
)

func NewUUID() string {
	v7, err := uuid.NewV7()
	if err != nil {
		slog.Error("failed to generate UUID v7, falling back to v4", slog.String("error", err.Error()))
		return uuid.New().String()
	}
	return v7.String()
}

func IsValid(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
