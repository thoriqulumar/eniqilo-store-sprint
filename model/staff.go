package model

import "time"

type Staff struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RegisterStaffRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,phone_number"`
	Name        string `json:"name" validate:"required,min=5,max=50"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type StaffServiceResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"accessToken"`
}

type RegisterStaffResponse struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
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
