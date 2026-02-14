package api_errors

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeValidation          string = "VALIDATION"
	CodeConflict            string = "CONFLICT"
	CodeInternalServerError string = "INTERNAL_SERVER_ERROR"
	CodeNotFound            string = "NOT_FOUND"
	CodeInvalidEnum         string = "INVALID_ENUM"
	CodeUnknown             string = "UNKNOWN"

	MessageDuplicateKey string = "Resource already exists"
	MessageNotFound     string = "Resource not found"
	MessageUnknown      string = "Something went wrong"
)

type ServiceError struct {
	Code    string
	Message string
}

func NewServiceError(code, message string) *ServiceError {
	return &ServiceError{Code: code, Message: message}
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func NewValidationError(message string) *ServiceError {
	return NewServiceError(CodeValidation, message)
}

func NewInternalServerError(message string) *ServiceError {
	return NewServiceError(CodeInternalServerError, message)
}

func FromDBError(err error) *ServiceError {
	if errors.Is(err, pgx.ErrNoRows) {
		return NewServiceError(CodeNotFound, "No rows found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return NewServiceError(CodeConflict, MessageDuplicateKey)
		case "23514":
			return NewServiceError(CodeInvalidEnum, pgErr.Message)
		case "23503":
			return NewServiceError(CodeNotFound, MessageNotFound)
		default:
			return NewServiceError(CodeUnknown, pgErr.Message)
		}
	}

	return NewServiceError(CodeUnknown, "Unknown error")
}
