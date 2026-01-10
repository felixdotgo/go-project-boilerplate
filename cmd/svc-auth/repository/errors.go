package repository

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrBadRequest            = errors.New("bad request")
	ErrOAuthProviderNotFound = errors.New("oauth provider not found")
	ErrProviderAlreadyLinked = errors.New("provider already linked to this account")
	ErrCannotUnlinkLastAuth  = errors.New("cannot unlink last authentication method")
	ErrRefreshTokenNotFound  = errors.New("refresh token not found")
	ErrRefreshTokenExpired   = errors.New("refresh token expired")
	ErrRefreshTokenRevoked   = errors.New("refresh token revoked")
)
