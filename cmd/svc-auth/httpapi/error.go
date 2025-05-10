package httpapi

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(code int, message string) error {
	return status.Error(toGrpcErrorCode(code), message)
}

func toGrpcErrorCode(code int) codes.Code {
	// convert HTTP status code to gRPC status code
	switch code {
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusConflict:
		return codes.AlreadyExists
	}
	return codes.Unknown
}
