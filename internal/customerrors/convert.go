package customerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConvertError converts an error into a gRPC error with the appropriate status code.
func ConvertError(err error) error {
	switch {
	case errors.Is(err, ErrNotFound):
		return status.Errorf(codes.NotFound, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
