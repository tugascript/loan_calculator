package api_errors

import (
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
