package types

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntityModel struct {
	Id        uint           `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
