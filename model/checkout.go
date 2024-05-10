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
	CustomerId     string          `json:"customerId"`
	ProductDetails []ProductDetail `json:"productDetails"`
	Paid           int             `json:"paid"`
	Change         int             `json:"change"`
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
