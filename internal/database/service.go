package database

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/utils"
)

var gormConfig = &gorm.Config{
	TranslateError: true,
}

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
		return nil, fmt.Errorf("unsupported database type, supported types are: mysql, postgres")
	}

	if err != nil {
		return nil, err
	}

	slog.Info("Database Connected Successfully")
	return database, nil
}

func CloseConnection(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	slog.Info("Database connection closed successfully")
	return nil
}

func ConnectSQL() (*gorm.DB, error) {
	dsn, err := GetDatabaseDSN()

	if err != nil {
		return nil, fmt.Errorf("failed to get SQL database DSN: %w", err)
	}

	database, err := gorm.Open(mysql.Open(dsn), gormConfig)

	if err != nil {
		return nil, fmt.Errorf("SQL database connection failed: %w", err)
	}

	return database, nil
}

func ConnectPostgres() (*gorm.DB, error) {
	dsn, err := GetDatabaseDSN()

	if err != nil {
		return nil, fmt.Errorf("failed to get Postgres database DSN: %w", err)
	}

	database, err := gorm.Open(postgres.Open(dsn), gormConfig)

	if err != nil {
		return nil, fmt.Errorf("postgres database connection failed: %w", err)
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
		mysqlConfig := mysqldriver.NewConfig()
		mysqlConfig.User = databaseUser
		mysqlConfig.Passwd = databasePassword
		mysqlConfig.Net = "tcp"
		mysqlConfig.Addr = net.JoinHostPort(databaseHost, databasePort)
		mysqlConfig.DBName = databaseName
		mysqlConfig.ParseTime = true
		return mysqlConfig.FormatDSN(), nil
	case Postgres:
		postgresURL := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(databaseUser, databasePassword),
			Host:     net.JoinHostPort(databaseHost, databasePort),
			Path:     "/" + databaseName,
			RawQuery: "sslmode=disable",
		}
		return postgresURL.String(), nil
	default:
		return "", errors.New("unsupported database type, supported types are: mysql, postgres")
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
		return fmt.Errorf("error during migration: %w", migrateError)
	}

	slog.Info("Migrations completed.")

	return nil

}
