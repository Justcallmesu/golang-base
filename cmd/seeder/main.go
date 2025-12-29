package main

import (
	"justcallmesu.com/golang-base/internal/app/users"
	"justcallmesu.com/golang-base/internal/config"
	"justcallmesu.com/golang-base/internal/database"
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
