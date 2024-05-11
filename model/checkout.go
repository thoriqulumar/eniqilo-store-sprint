package model

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	UserId      string    `json:"userId" db:"userId"`
	PhoneNumber string    `json:"phoneNumber" db:"phoneNumber"`
	Name        string    `json:"name" db:"name"`
	CreatedAt   time.Time `json:"-" db:"createdAt"`
}

type CustomerRequest struct {
	PhoneNumber *string `json:"phoneNumber" validate:"required,phone_number"`
	Name        *string `json:"name" validate:"required,min=5,max=50"`
}

type CustomerResponse struct {
	Message string   `json:"message"`
	Data    Customer `json:"data"`
}

type ErrorMessageOrder struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}

type ProductDetail struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderRequest struct {
	CustomerId     *string         `json:"customerId" validate:"required"`
	ProductDetails []ProductDetail `json:"productDetails"`
	Paid           *int            `json:"paid" validate:"required"`
	Change         *int            `json:"change" validate:"required"`
}

type CustomerResponseData struct {
	UserId      string `json:"userId" db:"userId"`
	PhoneNumber string `json:"phoneNumber" db:"phoneNumber"`
	Name        string `json:"name" db:"name"`
}
type ResponseCustomerList struct {
	Message string                 `json:"message"`
	Data    []CustomerResponseData `json:"data"`
}

type Transaction struct {
	TransactionId  uuid.UUID       `json:"transactionId" db:"transactionId"`
	CustomerId     uuid.UUID       `json:"customerId" db:"customerId"`
	ProductDetails []ProductDetail `json:"productDetails" db:"productDetails"`
	Paid           int             `json:"paid" db:"paid"`
	Change         int             `json:"change" db:"change"`
	CreatedAt      time.Time       `json:"createdAt" db:"createdAt"`
}

type GenericResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type GetHistoryParam struct {
	CustomerId *uuid.UUID
	Limit      int
	Offset     int
	CreatedAt  *string
}
