package main

import (
	"avitoTT/internal/repository"
	"avitoTT/internal/service"
	"avitoTT/openapi/handlers"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	db := repository.SetupDB()
	defer db.Close()

	serviceContainer := service.NewServiceContainer(db)
	c, err := handlers.NewContainer(*serviceContainer)
	if err != nil {
		e.Logger.Fatal("Error with initialize new container: ", err)
	}

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.POST("/api/auth", c.ApiAuthPost)
	e.GET("/api/buy/:item", c.ApiBuyItemGet)
	e.GET("/api/info", c.ApiInfoGet)
	e.POST("/api/sendCoin", c.ApiSendCoinPost)

	e.Logger.Fatal(e.Start(":8080"))
}
