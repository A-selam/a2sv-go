package infrastructure

import (
	domain "task_manager/Domain"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordService struct{}

func NewPasswordService() domain.PasswordService {
	return &BcryptPasswordService{}
}

func (s *BcryptPasswordService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *BcryptPasswordService) ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
