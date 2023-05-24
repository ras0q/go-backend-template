package config

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func getEnv(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return v
}

func AppAddr() string {
	return getEnv("APP_ADDR", ":8080")
}

func MySQL() *mysql.Config {
	return &mysql.Config{
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
}
