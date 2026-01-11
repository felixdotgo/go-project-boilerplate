package service

import "errors"

var (
	ErrBadRequest = errors.New("bad request")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrProviderAlreadyLinked = errors.New("provider already linked")

	ErrRefreshTokenInvalid  = errors.New("refresh token invalid")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")

	ErrOAuthExchangeFailed = errors.New("oauth exchange failed")
	ErrOAuthUserInfoFailed = errors.New("oauth user info failed")

	ErrOAuthCodeEmpty = errors.New("oauth code empty")

	ErrOAuthProviderNotConfigured = errors.New("oauth provider not configured")
	ErrOAuthProviderAlreadyLinked = errors.New("oauth provider already linked")
	ErrOAuthProviderEmpty         = errors.New("oauth provider empty")
	ErrOAuthProviderNotLinked     = errors.New("oauth provider not linked")

	ErrOAuthStateEmpty = errors.New("oauth state empty")

	ErrCannotUnlinkLastAuthMethod = errors.New("cannot unlink last authentication method")
	ErrCannotUnlinkLastAuth       = errors.New("cannot unlink last authentication method")
)
