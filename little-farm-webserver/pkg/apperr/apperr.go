package apperr

import (
	"errors"
	"fmt"
)

var (
	ErrConflict            = errors.New("conflict")
	ErrNotFound            = errors.New("not found")
	ErrUnprocessableEntity = errors.New("unprocessable entity")

	ErrValidation           = errors.New("validation error")
	ErrInvalidUUID          = fmt.Errorf("%w: id must be a valid UUID", ErrValidation)
	ErrMissingRequiredParam = fmt.Errorf("%w: missing required parameter", ErrValidation)

	ErrInternal         = errors.New("internal server error")
	ErrUnknownEventType = fmt.Errorf("%w: unknown event type", ErrInternal)
)

func IsPermanentError(err error) bool {
	return errors.Is(err, ErrValidation) || errors.Is(err, ErrUnknownEventType)
}
