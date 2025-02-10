package handlers

import (
	"avitoTT/openapi/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ApiAuthPost - Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
func (c *Container) ApiAuthPost(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld{
		Message: "Hello World",
	})
}

// ApiBuyItemGet - Купить предмет за монеты.
func (c *Container) ApiBuyItemGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.HelloWorld{
		Message: "Hello World",
	})
}

// ApiInfoGet - Получить информацию о монетах, инвентаре и истории транзакций.
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
