package crypto

import (
	"eniqilo-store/config"
	"eniqilo-store/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id, phoneNumber, name string) (string, error) {
	secret := config.GetString("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Id:          id,
		PhoneNumber: phoneNumber,
		Name:        name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
