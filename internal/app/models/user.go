package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"justcallmesu.com/golang-base/internal/common/types"
)

type User struct {
	types.BaseSoftDeleteModel

	// Identity & Contact
	Username string `json:"username,omitempty" gorm:"unique;not null"`
	Email    string `json:"email,omitempty" gorm:"unique;index"`
	FullName string `json:"fullName,omitempty" gorm:"index"`

	// Security & Authentication
	Password    string     `json:"-" gorm:"not null"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`

	// Relationships
	RoleId *types.UUID `json:"roleId,omitempty,omitzero"`
	Role   Role        `json:"role,omitempty,omitzero" gorm:"foreignKey:RoleId;references:ID;constraint:OnDelete:SET NULL"`
}

func (user *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	encryptedBytes, encryptionError := bcrypt.GenerateFromPassword([]byte(password), 10)

	if encryptionError != nil {
		return encryptionError
	}

	user.Password = string(encryptedBytes)
	return nil
}

func (user *User) ComparePassword(candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(candidatePassword))
}
