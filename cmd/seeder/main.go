package main

import (
	"justcallmesu.com/rest-api/internal/app/users"
	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
)

func main() {

	config.LoadConfig()

	database := database.InitConnection()

	users := []*users.User{
		{
			Username: "",
			Password: "",
		},
	}

	for _, value := range users {
		hashError := value.HashPassword()

		if hashError != nil {
			panic("Error hashing password: " + hashError.Error())
		}
	}

	database.Create(users)
}
