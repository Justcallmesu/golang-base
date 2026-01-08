package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/utils/default_value"
)

func InitConnection() *gorm.DB {
	databaseType := DatabaseType(default_value.GetDefaultValue(os.Getenv("DB_TYPE"), "mysql"))

	var database *gorm.DB

	switch databaseType {
	case MySQL:
		database = ConnectSQL()
	case Postgres:
		database = ConnectPostgres()
	default:
		panic("Unsupported database type, supported types are: mysql, postgres")
	}

	fmt.Println("Database Connected Successfully")
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
	databaseUser := default_value.GetDefaultValue(os.Getenv("DB_USER"), "root")
	databaseHost := default_value.GetDefaultValue(os.Getenv("DB_HOST"), "localhost")
	databasePort := default_value.GetDefaultValue(os.Getenv("DB_PORT"), "3306")
	databaseName := default_value.GetDefaultValue(os.Getenv("DB_NAME"), "golang_base")

	var databaseType = DatabaseType(default_value.GetDefaultValue(os.Getenv("DB_TYPE"), "mysql"))

	switch databaseType {
	case MySQL:
		databasePassword := default_value.GetDefaultValue(os.Getenv("DB_PASSWORD"), "")
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseName,
		)
	case Postgres:
		databasePassword := default_value.GetDefaultValue(os.Getenv("DB_PASSWORD"), "''")
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
