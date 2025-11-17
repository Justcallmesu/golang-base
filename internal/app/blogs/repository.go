package blogs

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BlogRepository struct {
	sqlDatabaseConnection *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{
		sqlDatabaseConnection: db,
	}
}

func (repository *BlogRepository) FindMany(parameter *BlogQuery, context context.Context) (*[]Blog, error) {

	var blogs []Blog

	qb := repository.sqlDatabaseConnection.WithContext(context)

	if parameter.Search != "" {
		searchValue := "%" + parameter.Search + "%"
		qb = qb.Where("Title LIKE ?", searchValue)
	}

	if parameter.OrderBy != "" && parameter.SortBy != "" {
		orderValue := fmt.Sprintf("%s %s", parameter.OrderBy, parameter.SortBy)
		qb = qb.Order(orderValue)
	}

	// Page is expected to be 1-based in QueryParams. Offset should be (page-1)*limit.
	offset := parameter.Limit * (parameter.Page - 1)
	qb = qb.Limit(parameter.Limit).Offset(offset)

	fetchError := qb.Find(&blogs).Error

	return &blogs, fetchError
}

func (repository *BlogRepository) FindOne(targetId int, context context.Context) (Blog, error) {

	var blog Blog

	fetchError := repository.sqlDatabaseConnection.WithContext(context).Preload("Details").Preload("Details.File").First(&blog, targetId).Error

	if fetchError != nil {
		return Blog{}, repository.HandleDatabaseError(fetchError)
	}

	return blog, nil
}

func (repository *BlogRepository) CreateOne(newBlog *Blog, context context.Context) error {
	result := gorm.WithResult()
	createError := gorm.G[Blog](repository.sqlDatabaseConnection, result).Create(context, newBlog)

	return createError
}

func (repository *BlogRepository) UpdateOne(updatedBlog *Blog, context context.Context) error {
	transaction := repository.sqlDatabaseConnection.WithContext(context).Begin()

	if transaction.Error != nil {
		return repository.HandleDatabaseError(transaction.Error)
	}

	defer func() {
		if transactionError := recover(); transactionError != nil {
			transaction.Rollback()
		}
	}()

	//  Update Parent Data
	if updateError := transaction.Model(updatedBlog).Omit("Details").Updates(updatedBlog).Error; updateError != nil {
		return repository.HandleDatabaseError(updateError)
	}

	// Update and Delete Association
	if updateDetailsError := transaction.Unscoped().Model(updatedBlog).Association("Details").Unscoped().Replace(updatedBlog.Details); updateDetailsError != nil {
		return repository.HandleDatabaseError(updateDetailsError)
	}

	return repository.HandleDatabaseError(transaction.Commit().Error)
}

func (repository *BlogRepository) DeleteOne(targetId int, context context.Context) error {
	transaction := repository.sqlDatabaseConnection.WithContext(context).Begin().Unscoped()

	defer func() {
		if transactionError := recover(); transactionError != nil {
			transaction.Rollback()
		}
	}()

	if deleteError := transaction.Where("id = ?", targetId).Delete(&Blog{}).Error; deleteError != nil {
		return repository.HandleDatabaseError(deleteError)
	}

	return repository.HandleDatabaseError(transaction.Commit().Error)
}

func (repository *BlogRepository) HandleDatabaseError(databaseError error) error {
	if databaseError == nil {
		return nil
	}

	if errors.Is(databaseError, gorm.ErrRecordNotFound) {
		return errors.New("Blog tidak ditemukan")
	}

	return errors.New("[Blog] Internal Server Error")
}
