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

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Errors: "Username and password are required",
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
	log.Println("Handlers: ApiBuyItemGet")

	items := ctx.Param("item")
	if items == "" {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Errors: "item parameter is required",
		})
	}

	token, err := c.BuyService.ExtractTokenFromHeader(ctx)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Errors: "Invalid or missing token",
		})
	}

	username, err := c.BuyService.ExtractUsernameFromToken(token)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Errors: "Invalid or missing token",
		})
	}

	err = c.BuyService.BuyItem(username, items)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Errors: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "Purchase successful"})
}

// ApiInfoGet - Получить информацию о монетах, инвентаре и истории транзаsкций.
func (c *Container) ApiInfoGet(ctx echo.Context) error {
	log.Println("Handlers: ApiInfoGet")

	return ctx.JSON(http.StatusOK, models.CurrentRequest{
		Message: "Hello World",
	})
}

// ApiSendCoinPost - Отправить монеты другому пользователю.
func (c *Container) ApiSendCoinPost(ctx echo.Context) error {
	log.Println("Handlers: ApiSendCoinPost")

	return ctx.JSON(http.StatusOK, models.CurrentRequest{
		Message: "Hello World",
	})
}
