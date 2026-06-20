package crypto

import (
	"frv-backend/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHashManager struct{}

func NewBcryptHashManager() ports.HashManager {
	return BcryptHashManager{}
}

func (h BcryptHashManager) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h BcryptHashManager) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
