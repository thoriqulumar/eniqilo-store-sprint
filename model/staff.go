package model

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	UserId      uuid.UUID `json:"userId" db:"userId"`
	Name        string    `json:"name" db:"name"`
	PhoneNumber string    `json:"phoneNumber" db:"phoneNumber"`
	Password    string    `json:"-" db:"password"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
}

type RegisterStaffRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,phone_number"`
	Name        string `json:"name" validate:"required,min=5,max=50"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type StaffWithToken struct {
	UserId      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	AccessToken string    `json:"accessToken"`
}

type LoginStaffRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,phone_number"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type RegisterStaffResponse struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	AccessToken string `json:"accessToken"`
}

func NewUser(phoneNumber, name, password string) *Staff {
	staff := &Staff{
		PhoneNumber: phoneNumber,
		Name:        name,
		Password:    password,
	}

	return staff
}
