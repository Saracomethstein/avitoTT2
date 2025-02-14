package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserCreationFailed = errors.New("failed to create user")
	ErrDatabaseIssue      = errors.New("database error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrMissingToken       = errors.New("authorization header is missing")
	ErrInvalidTokenFormat = errors.New("invalid authorization format")
	ErrNotEnoughCoins     = errors.New("not enough coins to buy item")
	ErrSendHimself        = errors.New("u can`t send coins to himself")
	ErrUserNotFound       = errors.New("user not found")
	ErrBalance            = errors.New("insufficient balance")
)
