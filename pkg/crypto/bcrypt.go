package crypto

import (
	"eniqilo-store/config"

	"golang.org/x/crypto/bcrypt"
)

var (
	cfg config.Config
)

func GenerateHashedPassword(password string, salt int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
