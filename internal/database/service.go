package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *sql.DB

func InitConnection() *gorm.DB {
	datasourceNetwork := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, databaseError := gorm.Open(mysql.New(mysql.Config{
		DSN: datasourceNetwork,
	}))

	if databaseError != nil {
		// This error is rare but should be handled.
		panic(fmt.Sprintf("Database connection error: %v", databaseError))

	}

	return database
}
