package handlers

import (
	"avitoTT/openapi/models"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ApiAuthPost - Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
func (c *Container) ApiAuthPost(ctx echo.Context) error {
	log.Println("Handlers: ApiAuthPost")
	var req models.AuthRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Errors: "Invalid request format: " + err.Error(),
		})
	}

	response, err := c.AuthService.Authenticate(req)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Errors: "Invalid username or password",
			})
		case errors.Is(err, models.ErrUserCreationFailed):
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Errors: "Failed to create user",
			})
		case errors.Is(err, models.ErrDatabaseIssue):
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Errors: "Database error",
			})
		default:
			log.Println("Unhandled error:", err)
			return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Errors: "Unknown server error",
			})
		}
	}

	return ctx.JSON(http.StatusOK, models.AuthResponse{
		Token: response.Token,
	})
}

// ApiBuyItemGet - Купить предмет за монеты.
func (c *Container) ApiBuyItemGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.CurrentRequest{
		Message: "Hello World",
	})
}

// ApiInfoGet - Получить информацию о монетах, инвентаре и истории транзаsкций.
func (c *Container) ApiInfoGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.CurrentRequest{
		Message: "Hello World",
	})
}

// ApiSendCoinPost - Отправить монеты другому пользователю.
func (c *Container) ApiSendCoinPost(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.CurrentRequest{
		Message: "Hello World",
	})
}
