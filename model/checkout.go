package model

import "time"

type Customer struct {
	UserId      string    `json:"userId" db:"userId"`
	PhoneNumber string    `json:"phoneNumber" db:"phoneNumber"`
	Name        string    `json:"name" db:"name"`
	CreatedAt   time.Time `json:"-" db:"createdAt"`
}

type CustomerRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}
