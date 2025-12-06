package integrationtests

import (
	"fmt"
	"log"
	"sync/atomic"
	"testing"

	"github.com/ras0q/go-backend-template/api"
	"github.com/ras0q/go-backend-template/core"
	"github.com/ras0q/go-backend-template/core/database"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

var globalServer atomic.Pointer[api.Server]

func TestMain(m *testing.M) {
	if err := run(m); err != nil {
		log.Fatalf("runtime error: %+v", err)
	}
}

func run(m *testing.M) error {
	config := core.Config{
		DBUser: "root",
		DBPass: "pass",
		DBHost: "localhost",
		DBPort: 3306,
		DBName: "app",
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("connect to docker: %w", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return fmt.Errorf("ping docker: %w", err)
	}

	mysqlConfig := config.MySQLConfig()

	resource, err := pool.Run("mysql", "latest", []string{
		"MYSQL_ROOT_PASSWORD=" + mysqlConfig.Passwd,
		"MYSQL_DATABASE=" + mysqlConfig.DBName,
	})
	if err != nil {
		return fmt.Errorf("start mysql docker: %w", err)
	}

	mysqlConfig.Addr = "localhost:" + resource.GetPort("3306/tcp")

	log.Println("wait for database container")

	var db *sqlx.DB
	if err := pool.Retry(func() error {
		_db, err := database.Setup(mysqlConfig)
		if err != nil {
			return err
		}

		db = _db

		return nil
	}); err != nil {
		return fmt.Errorf("connect to database container: %w", err)
	}

	deps := core.InjectDeps(db)

	server, err := api.NewServer(deps.Handler)
	globalServer.Store(server)

	m.Run()

	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("purge mysql docker: %w", err)
	}

	return nil
}
