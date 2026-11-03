package role_repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	role_dto "justcallmesu.com/golang-base/internal/app/dtos/role"
	"justcallmesu.com/golang-base/internal/app/models"
	"justcallmesu.com/golang-base/internal/common/apperror"
	"justcallmesu.com/golang-base/internal/common/types"
	"justcallmesu.com/golang-base/internal/common/utils"
)

type Repository struct {
	database *gorm.DB
}

var roleSortColumn = map[string]string{"name": "roles.name", "createdAt": "roles.created_at"}

func NewRepository(database *gorm.DB) *Repository {
	return &Repository{
		database: database,
	}
}

func (repository *Repository) FindAll(queryParams *role_dto.RoleQueryDto, context context.Context) (*types.PaginatedResponse[models.Role], error) {

	queryBuilder := repository.database.Model(&models.Role{}).Preload("Permissions", nil)

	if queryParams.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", queryParams.Search)
		queryBuilder = queryBuilder.Where("name LIKE ?", searchPattern)
	}

	if queryParams.SortBy != "" && queryParams.SortDir != "" {
		sortBy, sortDir, err := utils.ResolveSort(roleSortColumn, queryParams.SortBy, queryParams.SortDir)
		if err != nil {
			return nil, err
		}
		queryBuilder = queryBuilder.Order(fmt.Sprintf("%s %s", sortBy, sortDir))
	}

	paginatedResponse, fetchError := utils.Paginate[models.Role](queryBuilder, queryParams.Page, queryParams.PageSize)

	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}

	return paginatedResponse, nil
}

func (repository *Repository) FindByID(roleID types.UUID, context context.Context) (*models.Role, error) {
	role, fetchError := gorm.G[models.Role](repository.database).Preload("Permissions", nil).Where("ID = ?", roleID).First(context)

	if fetchError != nil {
		return nil, repository.handleDatabaseError(fetchError)
	}

	return &role, nil
}

func (repository *Repository) handleDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperror.NotFound("Role not found")

	}

	return apperror.Internal(err)
}
