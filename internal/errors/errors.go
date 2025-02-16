package handle_errors

import (
	"avitoTT/openapi/models"
	"errors"
	"net/http"
)

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

var errorMapping = map[error]int{
	ErrInvalidCredentials: http.StatusBadRequest,
	ErrUserCreationFailed: http.StatusInternalServerError,
	ErrDatabaseIssue:      http.StatusInternalServerError,
	ErrBalance:            http.StatusBadRequest,
	ErrSendHimself:        http.StatusBadRequest,
	ErrUserNotFound:       http.StatusBadRequest,
}

func Error(err error, defaultError string) (int, models.ErrorResponse) {
	for appErr, status := range errorMapping {
		if errors.As(err, &appErr) || errors.Is(err, appErr) {
			return status, models.ErrorResponse{Errors: appErr.Error()}
		}
	}

	return http.StatusInternalServerError, models.ErrorResponse{Errors: defaultError}
}
