package blogs

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BlogDetailsRepository struct {
	database *gorm.DB
}

func NewBlogDetailsRepository(database *gorm.DB) *BlogDetailsRepository {
	return &BlogDetailsRepository{
		database: database,
	}
}

func (repository BlogDetailsRepository) CreateOne(blogDetails *BlogDetails, context context.Context) error {
	result := gorm.WithResult()
	createError := gorm.G[BlogDetails](repository.database, result).Create(context, blogDetails)

	return createError
}

func (repository *BlogDetailsRepository) UpdateOne(blogDetails *BlogDetails, context context.Context) error {
	updateError := repository.database.Transaction(func(transaction *gorm.DB) error {

		updateError := transaction.Model(blogDetails).Updates(blogDetails).Error

		if updateError != nil {
			return updateError
		}

		return nil

	})

	if updateError != nil {
		return repository.HandleDatabaseError(updateError)
	}

	return nil
}

func (repository *BlogDetailsRepository) UpdateDetailOrder(updateBlogDetailsOrder *[]UpdateBlogDetailsOrder) error {
	updateError := repository.database.Transaction(func(tx *gorm.DB) error {
		for _, item := range *updateBlogDetailsOrder {
			if err := tx.Model(&BlogDetails{}).Where("id = ?", item.Id).UpdateColumn("order", item.Order).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if updateError != nil {
		return repository.HandleDatabaseError(updateError)
	}

	return nil
}

func (repository BlogDetailsRepository) HandleDatabaseError(databaseError error) error {
	if databaseError == nil {
		return nil
	}

	if errors.Is(databaseError, gorm.ErrRecordNotFound) {
		return errors.New("detail blog tidak ditemukan")
	}
	return errors.New("[Blog Detail] Internal Server Error")
}
