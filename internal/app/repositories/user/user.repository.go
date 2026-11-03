package user_repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	user_dtos "justcallmesu.com/golang-base/internal/app/dtos/user"
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type UserRepository struct {
	database *gorm.DB
}

var userSortColumn = map[string]string{
	"createdAt":   "users.created_at",
	"username":    "username",
	"fullName":    "full_name",
	"role":        "roles.name",
	"lastLoginAt": "last_login_at",
	"email":       "email",
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (repository *UserRepository) FindAll(queryParams *user_dtos.UserQueryDto, context context.Context) (*types.PaginatedResponse[models.User], error) {

	queryBuilder := repository.database.Model(&models.User{}).Preload("Role", nil).Joins("LEFT JOIN roles on roles.Id = users.role_id")

	if queryParams.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", queryParams.Search)
		queryBuilder = queryBuilder.Where("username LIKE ? OR roles.name LIKE ?", searchPattern, searchPattern)
	}

	if queryParams.RoleId != nil {
		queryBuilder = queryBuilder.Where("role_id = ?", *queryParams.RoleId)
	}

	if queryParams.SortBy != "" && queryParams.SortDir != "" {
		sortBy, sortDir, err := utils.ResolveSort(userSortColumn, queryParams.SortBy, queryParams.SortDir)

		if err != nil {
			return nil, err
		}

		queryBuilder = queryBuilder.Order(fmt.Sprintf("%s %s", sortBy, sortDir))
	}

	paginatedResponse, fetchError := utils.Paginate[models.User](queryBuilder, queryParams.Page, queryParams.PageSize)

	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}
	return paginatedResponse, nil

}

func (repository *UserRepository) FindUserByUsername(username string, context context.Context) (*models.User, error) {
	user, fetchError := gorm.G[models.User](repository.database).Where("Username = ?", username).Preload("Role", nil).Select("*").Preload("Role.Permissions", nil).First(context)
	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}

	return &user, nil
}

func (repository *UserRepository) FindUserByID(userID types.UUID, context context.Context) (*models.User, error) {
	user, fetchError := gorm.G[models.User](repository.database).Preload("Role", nil).Preload("Role.Permissions", nil).Where("ID = ?", userID).First(context)

	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}

	return &user, nil
}

func (repository *UserRepository) CreateUser(user *models.User, context context.Context) error {
	err := repository.database.Transaction(func(transaction *gorm.DB) error {
		return transaction.Save(&user).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *UserRepository) UpdateUser(user *models.User, context context.Context) error {
	err := repository.database.Transaction(func(transaction *gorm.DB) error {
		return transaction.Model(&user).Updates(&user).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *UserRepository) UpdateUserPassword(user *models.User, context context.Context) error {
	err := repository.database.Transaction(func(transaction *gorm.DB) error {
		return transaction.Save(&user).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *UserRepository) DeleteUser(userId types.UUID, context context.Context) error {
	err := repository.database.Transaction(func(transaction *gorm.DB) error {
		return transaction.Where("ID = ?", userId).Delete(&models.User{}).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *UserRepository) handleDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperror.NotFound("User not found")
	default:
		return apperror.Internal(err)
	}
}
