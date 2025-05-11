package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func Setup(mysqlConfig *mysql.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrateTables(db.DB); err != nil {
		return nil, err
	}

	return db, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func migrateTables(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("up migration: %w", err)
	}

	return nil
}
