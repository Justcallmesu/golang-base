package users

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	sqlDatabaseConnection *gorm.DB
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		sqlDatabaseConnection: database,
	}
}

func (repository *UserRepository) FindUserByUsername(email string, context context.Context) (*User, error) {
	user, fetchError := gorm.G[User](repository.sqlDatabaseConnection).First(context)

	if fetchError != nil {
		fmt.Printf("Error fetching user with email %s: %v\n", email, fetchError)

		if fetchError == gorm.ErrRecordNotFound {
			return nil, errors.New("User or Password Doesnt Match")
		}
		return nil, errors.New("something went wrong while fetching user")
	}

	return &user, nil
}
