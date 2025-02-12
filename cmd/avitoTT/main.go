package main

import (
	"avitoTT/internal/repository"
	"avitoTT/openapi/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	c, err := handlers.NewContainer()
	if err != nil {
		e.Logger.Fatal("Error with initialize new container: ", err)
	}

	db := repository.SetupDB()
	defer db.Close()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/auth", c.ApiAuthPost)
	e.GET("/api/buy/:item", c.ApiBuyItemGet)
	e.GET("/api/info", c.ApiInfoGet)
	e.POST("/api/sendCoin", c.ApiSendCoinPost)

	e.Logger.Fatal(e.Start(":8080"))
}
