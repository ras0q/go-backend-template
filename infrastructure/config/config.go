package config

import (
	"net"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/go-sql-driver/mysql"
)

type Config struct {
	AppAddr string `env:"APP_ADDR" default:":8080"`
	DBUser  string `env:"DB_USER" default:"root"`
	DBPass  string `env:"DB_PASS" default:"pass"`
	DBHost  string `env:"DB_HOST" default:"localhost"`
	DBPort  int    `env:"DB_PORT" default:"3306"`
	DBName  string `env:"DB_NAME" default:"app"`
}

func (c *Config) Parse() {
	kong.Parse(c)
}

func (c Config) MySQLConfig() *mysql.Config {
	mc := mysql.NewConfig()

	mc.User = c.DBUser
	mc.Passwd = c.DBPass
	mc.Net = "tcp"
	mc.Addr = net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort))
	mc.DBName = c.DBName
	mc.Collation = "utf8mb4_general_ci"
	mc.AllowNativePasswords = true

	return mc
}
