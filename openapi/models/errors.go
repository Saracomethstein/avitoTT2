package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserCreationFailed = errors.New("failed to create user")
	ErrDatabaseIssue      = errors.New("database error")
)
