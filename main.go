package main

import (
	"log"
	"net/http"

	"github.com/ras0q/go-backend-template/infrastructure/config"
	"github.com/ras0q/go-backend-template/infrastructure/database"
	"github.com/ras0q/go-backend-template/infrastructure/injector"
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

	var c config.Config
	c.Parse()

	// connect to and migrate database
	db, err := database.Setup(c.MySQLConfig())
	if err != nil {
		return err
	}
	defer g.Guard(db.Close)

	server, err := injector.InjectServer(db)
	if err != nil {
		return err
	}

	if err := http.ListenAndServe(c.AppAddr, server); err != nil {
		return err
	}

	return nil
}
