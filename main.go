package main

import (
	"log"

	"github.com/ras0q/go-backend-template/core"
	"github.com/ras0q/go-backend-template/core/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ras0q/goalie"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("runtime error: %+v", err)
	}
}

func run() (err error) {
	g := goalie.New()
	defer g.Collect(&err)

	var config core.Config
	config.Parse()

	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to and migrate database
	db, err := database.Setup(config.MySQLConfig())
	if err != nil {
		return err
	}
	defer g.Guard(db.Close)

	s := core.InjectDeps(db)

	core.SetupRoutes(s.Handler, e)

	if err := e.Start(config.AppAddr); err != nil {
		return err
	}

	return nil
}
