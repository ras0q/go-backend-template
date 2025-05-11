package main

import (
	"backend/cmd/server/injector"
	"backend/pkg/config"
	"backend/pkg/migration"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to database
	db, err := sqlx.Connect("mysql", config.MySQL().FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	// migrate tables
	if err := migration.MigrateTables(db.DB); err != nil {
		e.Logger.Fatal(err)
	}

	dep := injector.Inject(db)

	v1API := e.Group("/api/v1")
	dep.Handler.SetupRoutes(v1API)

	e.Logger.Fatal(e.Start(config.AppAddr()))
}
