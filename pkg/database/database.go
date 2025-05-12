package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
