package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func InitConnection() (*sql.DB, error) {
	datasourceNetwork := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", datasourceNetwork)

	if err != nil {
		// This error is rare but should be handled.
		return nil, err
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		// If ping fails, close the connection pool and return the error.
		dbCloseError := db.Close()

		if dbCloseError != nil {
			fmt.Printf("Failed to close database connection: %s\n", dbCloseError.Error())
		}

		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database Connected!")

	return db, nil
}
