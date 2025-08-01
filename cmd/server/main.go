package main

import (
	"backend/apps/server"
	"backend/pkg/config"
	"backend/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to and migrate database
	db, err := database.Setup(config.MySQL())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	s := server.InjectDeps(db)

	v1API := e.Group("/api/v1")
	server.SetupRoutes(s.Handler, v1API)

	e.Logger.Fatal(e.Start(config.AppAddr()))
}
