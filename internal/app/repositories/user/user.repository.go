package user_repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	user_dtos "justcallmesu.com/golang-base/internal/app/dtos/user"
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type Repository struct {
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

func NewRepository(database *gorm.DB) *Repository {
	return &Repository{
		database: database,
	}
}

func (repository *Repository) FindAll(queryParams *user_dtos.UserQueryDto, context context.Context) (*types.PaginatedResponse[models.User], error) {

	queryBuilder := repository.database.Model(&models.User{}).WithContext(context).Preload("Role", nil).Joins("LEFT JOIN roles on roles.Id = users.role_id")

	if queryParams.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", queryParams.Search)
		queryBuilder = queryBuilder.Where("users.username LIKE ? OR roles.name LIKE ?", searchPattern, searchPattern)
	}

	if queryParams.RoleId != nil {
		queryBuilder = queryBuilder.Where("users.role_id = ?", *queryParams.RoleId)
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

func (repository *Repository) FindUserByUsername(username string, context context.Context) (*models.User, error) {
	user, fetchError := gorm.G[models.User](repository.database).Where("Username = ?", username).Preload("Role", nil).Select("*").Preload("Role.Permissions", nil).First(context)
	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}

	return &user, nil
}

func (repository *Repository) FindUserByID(userID types.UUID, context context.Context) (*models.User, error) {
	user, fetchError := gorm.G[models.User](repository.database).Preload("Role", nil).Preload("Role.Permissions", nil).Where("ID = ?", userID).First(context)

	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}

	return &user, nil
}

func (repository *Repository) CreateUser(user *models.User, context context.Context) error {
	err := repository.database.WithContext(context).Transaction(func(transaction *gorm.DB) error {
		return transaction.Omit(clause.Associations).Save(&user).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *Repository) UpdateUser(user *models.User, context context.Context) error {
	err := repository.database.WithContext(context).Transaction(func(transaction *gorm.DB) error {
		return transaction.Omit(clause.Associations).Model(&user).Updates(&user).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *Repository) UpdateUserPassword(user *models.User, context context.Context) error {
	err := repository.database.WithContext(context).Transaction(func(transaction *gorm.DB) error {
		return transaction.Save(&user).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *Repository) DeleteUser(userId types.UUID, context context.Context) error {
	err := repository.database.WithContext(context).Transaction(func(transaction *gorm.DB) error {
		return transaction.Where("ID = ?", userId).Delete(&models.User{}).Error
	})

	return repository.handleDatabaseError(err)
}

func (repository *Repository) handleDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperror.NotFound("User not found")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return apperror.Conflict("user with the same username or email already exists")
	default:
		return apperror.Internal(err)
	}
}
