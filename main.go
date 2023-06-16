package main

import (
	"go-backend-sample/internal/handler"
	"go-backend-sample/internal/pkg/config"
	"go-backend-sample/internal/repository"

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

	// setup repository
	repo := repository.New(db)
	if err := repo.SetupTables(); err != nil {
		e.Logger.Fatal(err)
	}

	// setup routes
	h := handler.New(repo)
	v1API := e.Group("/api/v1")
	h.SetupRoutes(v1API)

	e.Logger.Fatal(e.Start(config.AppAddr()))
}
