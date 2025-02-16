package handle_errors_test

import (
	handle_errors "avitoTT/internal/errors"
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
			inputError:   handle_errors.ErrInvalidCredentials,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  handle_errors.ErrInvalidCredentials.Error(),
		},
		{
			name:         "User Creation Failed",
			inputError:   handle_errors.ErrUserCreationFailed,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  handle_errors.ErrUserCreationFailed.Error(),
		},
		{
			name:         "Database Issue",
			inputError:   handle_errors.ErrDatabaseIssue,
			defaultError: "internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  handle_errors.ErrDatabaseIssue.Error(),
		},
		{
			name:         "User Not Found",
			inputError:   handle_errors.ErrUserNotFound,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  handle_errors.ErrUserNotFound.Error(),
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
			inputError:   handle_errors.ErrDatabaseIssue,
			defaultError: "internal server error",
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  handle_errors.ErrDatabaseIssue.Error(),
		},
		{
			name:         "Wrapped Error (User Not Found)",
			inputError:   handle_errors.ErrUserNotFound,
			defaultError: "internal server error",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  handle_errors.ErrUserNotFound.Error(),
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
