package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BaseSoftDeleteModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID.UUID == uuid.Nil {
		newID, err := NewUUID()
		if err != nil {
			return err
		}
		m.ID = newID
	}
	return nil
}
