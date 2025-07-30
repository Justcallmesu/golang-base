package users

import (
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

func (repository *UserRepository) Create(userData *User) (*User, error) {

	return nil, nil
}

func (repository *UserRepository) FindUserByEmail(email string) (*User, error) {
	return nil, nil
}
