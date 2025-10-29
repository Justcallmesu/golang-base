package main

import (
	"fmt"

	"justcallmesu.com/rest-api/internal/app/blogs"
	"justcallmesu.com/rest-api/internal/app/users"
	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
)

func main() {

	config.LoadConfig()

	database := database.InitConnection()

	migrateError := database.AutoMigrate(
		&users.User{},
		&blogs.Blog{},
		&blogs.BlogDetails{},
	)

	if migrateError != nil {
		panic("Error migrating database: " + migrateError.Error())
	}

	fmt.Println("Migrations finished")

}
