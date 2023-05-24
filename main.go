package main

import (
	"fmt"
	"go-backend-sample/handler"
	"go-backend-sample/repository"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	appAddr  = getEnv("APP_ADDR", ":8080")
	dbConfig = mysql.Config{
		User:   getEnv("DB_USER", "root"),
		Passwd: getEnv("DB_PASSWORD", "pass"),
		Net:    getEnv("DB_NET", "tcp"),
		Addr: fmt.Sprintf(
			"%s:%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
		),
		DBName: getEnv("DB_NAME", "backend_sample"),
	}
)

func main() {
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// dependencies
	db, err := sqlx.Connect("mysql", dbConfig.FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
	}
	repo := repository.New(db)
	h := handler.New(repo)

	// routes
	api := e.Group("/api")
	h.SetupRoutes(api)

	e.Logger.Fatal(e.Start(appAddr))
}

func getEnv(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return v
}
