package database

import (
	"fmt"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/utils"
)

func InitConnection() (*gorm.DB, error) {
	databaseType := DatabaseType(utils.GetDefaultValue(os.Getenv("DB_TYPE"), "mysql"))

	var database *gorm.DB
	var err error

	switch databaseType {
	case MySQL:
		database, err = ConnectSQL()
	case Postgres:
		database, err = ConnectPostgres()
	default:
		return nil, fmt.Errorf("Unsupported database type, supported types are: mysql, postgres")
	}

	if err != nil {
		return nil, fmt.Errorf("Database connection failed: %v", err)
	}

	slog.Info("Database Connected Successfully")
	return database, nil
}

func ConnectSQL() (*gorm.DB, error) {
	dsn, err := GetDatabaseDSN()

	if err != nil {
		return nil, fmt.Errorf("Failed to get database DSN: %v", err)
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("SQL Database connection failed: %v", err)
	}

	return database, nil
}

func ConnectPostgres() (*gorm.DB, error) {
	dsn, err := GetDatabaseDSN()

	if err != nil {
		return nil, fmt.Errorf("Failed to get database DSN: %v", err)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Postgres Database connection failed: %v", err)
	}

	return database, nil
}

func GetDatabaseDSN() (string, error) {
	databaseUser := utils.GetDefaultValue(os.Getenv("DB_USER"), "root")
	databaseHost := utils.GetDefaultValue(os.Getenv("DB_HOST"), "localhost")
	databasePort := utils.GetDefaultValue(os.Getenv("DB_PORT"), "3306")
	databaseName := utils.GetDefaultValue(os.Getenv("DB_NAME"), "golang_base")

	databaseType := DatabaseType(utils.GetDefaultValue(os.Getenv("DB_TYPE"), "mysql"))
	databasePassword := utils.GetDefaultValue(os.Getenv("DB_PASSWORD"), "")

	switch databaseType {
	case MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseName,
		), nil
	case Postgres:
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			databaseHost,
			databasePort,
			databaseUser,
			databasePassword,
			databaseName,
		), nil
	default:
		return "", fmt.Errorf("Unsupported database type, supported types are: mysql, postgres")
	}
}

func Migrations(database *gorm.DB) error {

	slog.Warn("Running migrations...")

	migrateError := database.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
	)

	if migrateError != nil {
		return fmt.Errorf("Error during migration: %v", migrateError)
	}

	slog.Info("Migrations completed.")

	return nil

}
