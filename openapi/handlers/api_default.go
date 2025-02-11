package handlers

import (
	"avitoTT/openapi/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ApiAuthPost - Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
func (c *Container) ApiAuthPost(ctx echo.Context) error {
	var req models.AuthRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Errors: err.Error(),
		})
	}

	response, err := c.AuthService.Authenticate(req)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Errors: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, models.AuthResponse{
		Token: response.Message,
	})
}

// ApiBuyItemGet - Купить предмет за монеты.
func (c *Container) ApiBuyItemGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld{
		Message: "Hello World",
	})
}

// ApiInfoGet - Получить информацию о монетах, инвентаре и истории транзаsкций.
func (c *Container) ApiInfoGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld{
		Message: "Hello World",
	})
}

// ApiSendCoinPost - Отправить монеты другому пользователю.
func (c *Container) ApiSendCoinPost(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld{
		Message: "Hello World",
	})
}
