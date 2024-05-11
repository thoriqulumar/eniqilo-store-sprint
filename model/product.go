package model

import (
	"time"

	"github.com/google/uuid"
)

// Category defines product category type
type Category string

const (
	Clothing    Category = "Clothing"
	Accessories Category = "Accessories"
	Footwear    Category = "Footwear"
	Beverages   Category = "Beverages"
)

type GetProductParam struct {
	ID          *uuid.UUID
	Limit       *int
	Offset      *int
	Name        *string
	IsAvailable *bool
	Category    *Category
	SKU         *string
	InStock     *bool
	Sort        ProductSorting
}

type ProductSorting struct {
	Price     *string
	CreatedAt *string
}

type GetProductResponse struct {
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type PostProductResponse struct {
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type UpdateProductResponse struct {
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

// Product represents the product table structure
type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	SKU         string    `json:"sku" db:"sku"`
	Category    string    `json:"category" db:"category"`
	Stock       *int      `json:"stock" db:"stock"`
	Price       int       `json:"price" db:"price"`
	ImageURL    string    `json:"imageUrl" db:"imageUrl"`
	Notes       string    `json:"notes" db:"notes"`
	IsAvailable *bool     `json:"isAvailable" db:"isAvailable" `
	Location    string    `json:"location" db:"location"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
}

type Data struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}
