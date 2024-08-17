package customerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConvertError converts an error into a gRPC error with the appropriate status code.
func ConvertError(err error) error {
	// Check if the error is already a gRPC error
	if _, ok := status.FromError(err); ok {
		return err
	}

	var errNotFound *ErrNotFound

	switch {
	case errors.As(err, &errNotFound):
		return status.Errorf(codes.NotFound, errNotFound.Error())

	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
