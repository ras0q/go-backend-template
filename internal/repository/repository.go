package repository

import (
	"fmt"

	_ "embed"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

//go:embed schema.sql
var schema string

func (r *Repository) SetupTables() error {
	if _, err := r.db.Exec(schema); err != nil {
		return fmt.Errorf("setup tables: %w", err)
	}

	return nil
}
