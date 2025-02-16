package handle_errors_test

import (
	handle_errors "avitoTT/internal/errors"
	"avitoTT/openapi/models"
	"errors"
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		defaultError string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Invalid Credentials",
			inputError:   models.ErrInvalidCredentials,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  models.ErrInvalidCredentials.Error(),
		},
		{
			name:         "User Creation Failed",
			inputError:   models.ErrUserCreationFailed,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  models.ErrUserCreationFailed.Error(),
		},
		{
			name:         "Database Issue",
			inputError:   models.ErrDatabaseIssue,
			defaultError: "internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  models.ErrDatabaseIssue.Error(),
		},
		{
			name:         "User Not Found",
			inputError:   models.ErrUserNotFound,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  models.ErrUserNotFound.Error(),
		},
		{
			name:         "Unknown Error",
			inputError:   errors.New("some unknown error"),
			defaultError: "internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "internal server error",
		},
		{
			name:         "Wrapped Error (Database Issue)",
			inputError:   models.ErrDatabaseIssue,
			defaultError: "internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  models.ErrDatabaseIssue.Error(),
		},
		{
			name:         "Wrapped Error (User Not Found)",
			inputError:   models.ErrUserNotFound,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  models.ErrUserNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, errResp := handle_errors.Error(tt.inputError, tt.defaultError)

			if status != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, status)
			}
			if errResp.Errors != tt.expectedMsg {
				t.Errorf("expected error message %q, got %q", tt.expectedMsg, errResp.Errors)
			}
		})
	}
}
