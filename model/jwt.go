package model

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

type JWTPayload struct {
	Id          string
	Name        string
	PhoneNumber string
}
