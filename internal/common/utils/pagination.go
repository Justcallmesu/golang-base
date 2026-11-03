package utils

import (
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/common/types"
)

func Paginate[T any](database *gorm.DB, page int, pageSize int) (*types.PaginatedResponse[T], error) {

	var totalItems int64
	var items []T

	databaseModel := database.Model(new(T))
	offset := (page - 1) * pageSize

	countError := databaseModel.Session(&gorm.Session{}).Count(&totalItems).Error

	if countError != nil {
		return nil, countError
	}

	fetchError := databaseModel.Session(&gorm.Session{}).Limit(pageSize).Offset(offset).Find(&items).Error

	if fetchError != nil {
		return nil, fetchError
	}

	return &types.PaginatedResponse[T]{
		Data: items,
		Meta: types.PaginationMetadata{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  int((totalItems + int64(pageSize) - 1) / int64(pageSize)),
			TotalItems:  int(totalItems),
		},
	}, nil
}
