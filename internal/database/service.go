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

func InitConnection() *gorm.DB {
	databaseType := DatabaseType(utils.GetDefaultValue(os.Getenv("DB_TYPE"), "mysql"))

	var database *gorm.DB

	switch databaseType {
	case MySQL:
		database = ConnectSQL()
	case Postgres:
		database = ConnectPostgres()
	default:
		panic("Unsupported database type, supported types are: mysql, postgres")
	}

	slog.Info("Database Connected Successfully")
	return database
}

func ConnectSQL() *gorm.DB {
	dsn := GetDatabaseDSN()

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("SQL Database connection failed: %v", err))
	}

	return database
}

func ConnectPostgres() *gorm.DB {
	dsn := GetDatabaseDSN()

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Postgres Database connection failed: %v", err))
	}

	return database
}

func GetDatabaseDSN() string {
	databaseUser := utils.GetDefaultValue(os.Getenv("DB_USER"), "root")
	databaseHost := utils.GetDefaultValue(os.Getenv("DB_HOST"), "localhost")
	databasePort := utils.GetDefaultValue(os.Getenv("DB_PORT"), "3306")
	databaseName := utils.GetDefaultValue(os.Getenv("DB_NAME"), "golang_base")

	var databaseType = DatabaseType(utils.GetDefaultValue(os.Getenv("DB_TYPE"), "mysql"))

	switch databaseType {
	case MySQL:
		databasePassword := utils.GetDefaultValue(os.Getenv("DB_PASSWORD"), "")
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseName,
		)
	case Postgres:
		databasePassword := utils.GetDefaultValue(os.Getenv("DB_PASSWORD"), "''")
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			databaseHost,
			databasePort,
			databaseUser,
			databasePassword,
			databaseName,
		)
	default:
		panic("Unsupported database type, supported types are: mysql, postgres")
	}
}

func Migrations(database *gorm.DB) {

	slog.Warn("Running migrations...")

	migrateError := database.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
	)

	if migrateError != nil {
		panic("Error migrating database: " + migrateError.Error())
	}

	slog.Info("Migrations completed.")

}
