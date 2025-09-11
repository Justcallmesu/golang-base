package main

import (
	"context"
	"os/user"

	"gorm.io/gorm"
	"justcallmesu.com/rest-api/internal/app/users"
	"justcallmesu.com/rest-api/internal/config"
	"justcallmesu.com/rest-api/internal/database"
)

func main() {

	config.LoadConfig()
	ctx := context.Background()

	database := database.InitConnection()
	
	users := []*users.User{
		{
			Model: gorm.Model{
				ID:1,
			},
			Username: "justcallmesu",
			Password: "portfolioAppPasswordForJustcallmesu12361239789$!@#!@#!@*#@!)!",
		},
	}

	
	for _, value := range users {
		hashError := value.HashPassword()
		_,deleteError := gorm.G[user.User](database).Where("id = ?", value.ID).Delete(ctx)

		if(deleteError != nil){

			switch deleteError.Error() {
			case gorm.ErrRecordNotFound.Error():
				// Record not found, nothing to delete
			default:
				panic("Error deleting existing user: " + deleteError.Error())
			}

		}

		if hashError != nil {
			panic("Error hashing password: " + hashError.Error())
		}
	}

	database.Create(users)
}
