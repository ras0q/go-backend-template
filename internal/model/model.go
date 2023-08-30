package model

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() error {
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 3306
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "portal"
	}

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%d)/gorm?charset=utf8&parseTime=True&loc=Local", user, password, host, port), // data source name
		DefaultStringSize:         256,                                                                                                    // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                   // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                   // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                   // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                  // auto configure based on currently MySQL version
	}), &gorm.Config{})

	return nil

}
