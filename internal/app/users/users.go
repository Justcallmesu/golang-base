package users

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

func (user *User) HashPassword() error {

	encryptedBytes, encryptionError := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if encryptionError != nil {
		return encryptionError
	}

	user.Password = string(encryptedBytes)
	return nil
}

func (user *User) ComparePassword(candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(candidatePassword))
}
