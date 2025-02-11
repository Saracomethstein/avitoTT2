package main

import (
	"avitoTT/openapi/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	//todo: handle the error!
	c, err := handlers.NewContainer()
	if err != nil {
		e.Logger.Fatal("Error with initialize new container: ", err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ApiAuthPost - Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	e.POST("/api/auth", c.ApiAuthPost)

	// ApiBuyItemGet - Купить предмет за монеты.
	e.GET("/api/buy/:item", c.ApiBuyItemGet)

	// ApiInfoGet - Получить информацию о монетах, инвентаре и истории транзакций.
	e.GET("/api/info", c.ApiInfoGet)

	// ApiSendCoinPost - Отправить монеты другому пользователю.
	e.POST("/api/sendCoin", c.ApiSendCoinPost)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
