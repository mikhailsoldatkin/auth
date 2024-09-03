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
	var errInvalidPassword *ErrInvalidPassword
	var errInvalidToken *ErrInvalidToken
	var errForbidden *ErrForbidden

	switch {
	case errors.As(err, &errNotFound):
		return status.Errorf(codes.NotFound, errNotFound.Error())
	case errors.As(err, &errInvalidPassword):
		return status.Errorf(codes.Unauthenticated, errInvalidPassword.Error())
	case errors.As(err, &errInvalidToken):
		return status.Errorf(codes.Unauthenticated, errInvalidToken.Error())
	case errors.As(err, &errForbidden):
		return status.Errorf(codes.PermissionDenied, errForbidden.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
