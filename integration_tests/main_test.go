package integration_tests

import (
	"backend/cmd/server/server"
	"backend/pkg/config"
	"backend/pkg/database"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ory/dockertest/v3"
)

var e *echo.Echo

func TestMain(m *testing.M) {
	e = echo.New()
	e.Logger.SetLevel(log.INFO)

	pool, err := dockertest.NewPool("")
	if err != nil {
		e.Logger.Fatalf("connect to docker: %v", err)
	}

	if err := pool.Client.Ping(); err != nil {
		e.Logger.Fatalf("ping docker: %v", err)
	}

	mysqlConfig := config.MySQL()

	resource, err := pool.Run("mysql", "latest", []string{
		"MYSQL_ROOT_PASSWORD=" + mysqlConfig.Passwd,
		"MYSQL_DATABASE=" + mysqlConfig.DBName,
	})
	if err != nil {
		e.Logger.Fatalf("run docker: %v", err)
	}

	mysqlConfig.Addr = "localhost:" + resource.GetPort("3306/tcp")

	e.Logger.Info("wait for database container")

	var db *sqlx.DB
	if err := pool.Retry(func() error {
		_db, err := database.Setup(mysqlConfig)
		if err != nil {
			return err
		}

		db = _db

		return nil
	}); err != nil {
		e.Logger.Fatalf("connect to database container: %v", err)
	}

	s := server.Inject(db)
	s.SetupRoutes(e.Group("/api/v1"))

	e.Logger.Info("start integration test")
	m.Run()

	if err := pool.Purge(resource); err != nil {
		e.Logger.Fatalf("purge docker: %v", err)
	}
}
