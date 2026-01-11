package httpapi

import (
	"errors"
	"net/http"

	"github.com/0x46656C6978/go-project-boilerplate/cmd/svc-auth/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewError(code int, message string) error {
	return status.Error(toGrpcErrorCode(code), message)
}

// toGrpcErrorCode converts HTTP status code to gRPC status code
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

// FromUserToHTTPError maps service layer errors to HTTP errors
func FromUserToHTTPError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, service.ErrOAuthProviderNotConfigured) {
		return NewError(http.StatusBadRequest, "provider not configured")
	}
	if errors.Is(err, service.ErrOAuthProviderEmpty) {
		return NewError(http.StatusBadRequest, "provider is required")
	}
	if errors.Is(err, service.ErrOAuthStateEmpty) {
		return NewError(http.StatusBadRequest, "state is required")
	}
	if errors.Is(err, service.ErrOAuthCodeEmpty) {
		return NewError(http.StatusBadRequest, "code is required")
	}
	if errors.Is(err, service.ErrOAuthExchangeFailed) {
		return NewError(http.StatusBadRequest, "failed to exchange authorization code")
	}
	if errors.Is(err, service.ErrOAuthUserInfoFailed) {
		return NewError(http.StatusBadRequest, "failed to get user info from provider")
	}
	if errors.Is(err, service.ErrUserNotFound) {
		return NewError(http.StatusNotFound, "user not found")
	}
	if errors.Is(err, service.ErrOAuthProviderAlreadyLinked) {
		return NewError(http.StatusConflict, "provider already linked to another account")
	}
	if errors.Is(err, service.ErrOAuthProviderNotLinked) {
		return NewError(http.StatusNotFound, "provider not linked")
	}
	if errors.Is(err, service.ErrCannotUnlinkLastAuthMethod) {
		return NewError(http.StatusBadRequest, "cannot unlink last authentication method")
	}
	return NewError(http.StatusInternalServerError, ErrInternalServerError)
}
