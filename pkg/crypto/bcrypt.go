package crypto

import (
	"eniqilo-store/config"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(password string) (string, error) {
	costStr := config.GetString("BCRYPT_SALT")

	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
