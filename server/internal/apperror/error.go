package apperror

import (
	"encoding/json"
)

var (
	ErrNotFound         = NewAppError(nil, "not found", "BS-000001")
	ErrBadRequest       = NewAppError(nil, "bad request", "BS-000002")
	ErrInternalError    = NewAppError(nil, "internal error", "BS-000000")
	ErrUnauthorized     = NewAppError(nil, "unauthorized", "BS-000003")
	ErrPermissionDenied = NewAppError(nil, "permission denied", "BS-000004")
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
