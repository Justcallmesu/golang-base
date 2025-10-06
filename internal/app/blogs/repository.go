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

	queryBuilder := gorm.G[Blog](repository.sqlDatabaseConnection)

	if parameter.Search != "" {
		searchValue := "%" + parameter.Search + "%"
		queryBuilder.Where("Title Like ?", searchValue)
	}

	if parameter.OrderBy != "" && parameter.SortBy != "" {
		orderValue := fmt.Sprintf("%s %s", parameter.OrderBy, parameter.SortBy)
		queryBuilder.Order(orderValue)
	}

	queryBuilder.Limit(parameter.Limit).Offset(parameter.Limit * parameter.Page)

	blogs, fetchError := queryBuilder.Find(context)

	return &blogs, fetchError
}

func (repository *BlogRepository) FindOne(targetId int, context context.Context) (Blog, error) {
	foundBlog, fetchOneError := gorm.G[Blog](repository.sqlDatabaseConnection).Where("id = ?", targetId).Take(context)

	if fetchOneError != nil {
		return Blog{}, repository.HandleDatabaseError(fetchOneError)
	}

	return foundBlog, (repository.HandleDatabaseError(fetchOneError))
}

func (repository *BlogRepository) CreateOne(newBlog *Blog, context context.Context) error {
	result := gorm.WithResult()
	createError := gorm.G[Blog](repository.sqlDatabaseConnection, result).Create(context, newBlog)

	return createError
}

func (repository *BlogRepository) UpdateOne(updatedBlog *Blog, context context.Context) (int, error) {
	return gorm.G[Blog](repository.sqlDatabaseConnection).Where("id = ?", updatedBlog.ID).Updates(context, *updatedBlog)
}

func (repository *BlogRepository) DeleteOne(targetId int, context context.Context) (int, error) {
	return gorm.G[Blog](repository.sqlDatabaseConnection).Where("id = ?", targetId).Delete(context)

}

func (repositorry *BlogRepository) HandleDatabaseError(databaseError error) error {
	if databaseError == nil {
		return nil
	}

	if errors.Is(databaseError, gorm.ErrRecordNotFound) {
		return errors.New("Blog tidak ditemukan")
	}

	return errors.New("[Blog] Internal Server Error")
}
