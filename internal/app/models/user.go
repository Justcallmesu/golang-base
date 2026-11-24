package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"justcallmesu.com/golang-base/internal/common/types"
)

const passwordHashCost = 10

var dummyPasswordHash, _ = bcrypt.GenerateFromPassword([]byte("asdgasdjkasada"), passwordHashCost)

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

	encryptedBytes, encryptionError := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)

	if encryptionError != nil {
		return encryptionError
	}

	user.Password = string(encryptedBytes)
	return nil
}

func (user *User) ComparePassword(candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(candidatePassword))
}

// CompareDummyPassword takes as long as ComparePassword but always fails.
// Call it when the user does not exist so response time can't reveal which usernames exist.
func CompareDummyPassword(candidatePassword string) {
	_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(candidatePassword))
}
