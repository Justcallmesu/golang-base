package files

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type FilesRepository struct {
	database *gorm.DB
}

func NewFilesRepository(database *gorm.DB) *FilesRepository {
	return &FilesRepository{
		database: database,
	}
}

func (repository FilesRepository) FindAll(query FilesQuery, context context.Context) (*[]Files, error) {
	var fetchedFiles []Files

	queryBuilder := repository.database.WithContext(context)

	if query.Search != "" {
		searchValue := "%" + query.Search + "%"

		queryBuilder.Where("originalName LIKE ?", searchValue)
	}

	if query.OrderBy != "" && query.SortBy != "" {
		orderValue := fmt.Sprintf("%s %s", query.OrderBy, query.SortBy)

		queryBuilder.Order(orderValue)
	}

	queryBuilder.Limit(query.Limit).Offset(query.Limit * query.Page)

	fetchError := queryBuilder.Find(&fetchedFiles).Error

	return &fetchedFiles, fetchError
}

func (repository FilesRepository) FindOne(id int, context context.Context) (*Files, error) {
	var file Files

	fetchError := repository.database.WithContext(context).First(&file, id).Error

	if fetchError != nil {
		return nil, repository.HandleDatabaseError(fetchError)
	}

	return &file, nil
}

func (repository FilesRepository) CreateOne(file *Files, context context.Context) error {
	result := gorm.WithResult()
	createError := gorm.G[Files](repository.database, result).Create(context, file)

	return createError
}

func (repository FilesRepository) CreateMany(files *[]Files, context context.Context) error {
	result := gorm.WithResult()

	createError := gorm.G[[]Files](repository.database, result).Create(context, files)

	return createError
}

func (repository FilesRepository) UpdateOne(updatedFile *Files, context context.Context) error {
	transaction := repository.database.WithContext(context).Begin()

	if transaction.Error != nil {
		return repository.HandleDatabaseError(transaction.Error)
	}

	defer func() {
		if transactionError := recover(); transactionError != nil {
			transaction.Rollback()
		}
	}()

	if updateError := transaction.Model(updatedFile).Updates(updatedFile).Error; updateError != nil {
		return repository.HandleDatabaseError(updateError)
	}

	return repository.HandleDatabaseError(transaction.Commit().Error)
}

func (repository FilesRepository) DeleteOne(id int, context context.Context) error {
	transaction := repository.database.WithContext(context).Begin().Unscoped()

	if transaction.Error != nil {
		return repository.HandleDatabaseError(transaction.Error)
	}

	defer func() {
		if transactionError := recover(); transactionError != nil {
			transaction.Rollback()
		}
	}()

	if deleteError := transaction.Where("id = ?", id).Delete(Files{}).Error; deleteError != nil {
		return repository.HandleDatabaseError(deleteError)
	}

	return nil
}

func (repository FilesRepository) HandleDatabaseError(databaseError error) error {
	if databaseError == nil {
		return nil
	}

	if errors.Is(databaseError, gorm.ErrRecordNotFound) {
		return errors.New("file tidak ditemukan")
	}

	return errors.New("[File] Internal Server Error")
}
