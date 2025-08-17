package main

import (
	"backend/core"
	"backend/core/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var config core.Config
	config.Parse()

	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to and migrate database
	db, err := database.Setup(config.MySQLConfig())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	s := core.InjectDeps(db)

	core.SetupRoutes(s.Handler, e)

	e.Logger.Fatal(e.Start(config.AppAddr))
}
