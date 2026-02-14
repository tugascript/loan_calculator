package api_errors

import (
	"errors"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCStatus(e *ServiceError) error {
	switch e.Code {
	case CodeValidation, CodeInvalidEnum:
		return status.Error(codes.InvalidArgument, e.Message)
	case CodeConflict:
		return status.Error(codes.AlreadyExists, e.Message)
	case CodeNotFound:
		return status.Error(codes.NotFound, e.Message)
	case CodeUnknown:
		return status.Error(codes.Internal, e.Message)
	default:
		return status.Error(codes.Internal, MessageUnknown)
	}
}

func toSnakeCase(camel string) string {
	if camel == strings.ToUpper(camel) {
		return strings.ToLower(camel)
	}

	var result strings.Builder
	for i, char := range camel {
		if unicode.IsUpper(char) {
			lowered := unicode.ToLower(char)
			if i > 0 {
				result.WriteRune('_')
				result.WriteRune(lowered)
				continue
			}

			result.WriteRune(lowered)
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func InvalidArgumentResponseFromErr(err error) error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return status.Error(codes.InvalidArgument, "Invalid request")
	}

	fields := make([]string, len(errs))

	for i, field := range errs {
		fields[i] = toSnakeCase(field.Field())
	}

	errMessage := strings.Join(fields, ", ")
	if len(fields) > 1 {
		errMessage += " are invalid"
	} else {
		errMessage += " is invalid"
	}

	return status.Error(codes.InvalidArgument, errMessage)
}
