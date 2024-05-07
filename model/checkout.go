package model

type Customer struct {
	UserId      string `json:"userId" db:"userId"`
	PhoneNumber string `json:"phoneNumber" db:"phoneNumber"`
	Name        string `json:"name" db:"name"`
}

type CustomerRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}