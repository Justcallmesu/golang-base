package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UUID struct{ uuid.UUID }

func (UUID) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "uuid"
	default:
		return "char(36)"
	}
}

func NewUUID() (UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return UUID{}, err
	}
	return UUID{UUID: id}, nil
}

func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID{UUID: id}, nil
}

func (u *UUID) UnmarshalParam(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	u.UUID = id
	return nil
}
