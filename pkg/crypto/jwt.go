package crypto

import (
	"eniqilo-store/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(id uuid.UUID, phoneNumber, name, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Id:          id.String(),
		PhoneNumber: phoneNumber,
		Name:        name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func VerifyToken(token, secretKey string) (*model.JWTPayload, error) {
	claims := &model.JWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, err
	}

	payload := &model.JWTPayload{
		Id:          claims.Id,
		PhoneNumber: claims.PhoneNumber,
		Name:        claims.Name,
	}

	return payload, nil
}
