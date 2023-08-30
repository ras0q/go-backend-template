package main

import (
	"go-backend-sample/internal/model"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to database
	err := model.Setup()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.HideBanner = true
	e.HidePort = true
	e.Debug = true
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
