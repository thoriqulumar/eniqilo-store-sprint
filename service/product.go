package service

import (
	"context"
	"eniqilo-store/model"
	"eniqilo-store/repo"
	"errors"
	"net/url"

	"github.com/google/uuid"
)

// ProductService handles business logic related to products.
type ProductService interface {
	CreateProduct(ctx context.Context, data model.CreatedProduct) (model.CreatedProduct, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	repo repo.ProductRepo
}

// NewProductService creates a new instance of ProductService.
func NewProductService(repo repo.ProductRepo) ProductService {
	return &productService{
		repo: repo,
	}
}

// CreateProduct handles the creation of a new product.
func (s *productService) CreateProduct(ctx context.Context, prod model.CreatedProduct) (model.CreatedProduct, error) {
	// Validate the product
	if err := validateCreateProduct(prod); err != nil {
		return model.CreatedProduct{}, err
	}

	return s.repo.CreateProduct(ctx, prod)
}

func validateCreateProduct(prod model.CreatedProduct) error {
	// Name validation
	if prod.Name == "" || len(prod.Name) > 30 {
		return errors.New("name must not be empty and should be between 1 and 30 characters long")
	}

	// SKU validation
	if prod.SKU == "" || len(prod.SKU) > 30 {
		return errors.New("SKU must not be empty and should be between 1 and 30 characters long")
	}

	// Category validation
	validCategories := map[model.Category]bool{
		model.Clothing:    true,
		model.Accessories: true,
		model.Footwear:    true,
		model.Beverages:   true,
	}
	if _, ok := validCategories[model.Category(prod.Category)]; !ok {
		return errors.New("invalid category")
	}

	// Image URL validation
	// You can use regex or a library like net/url to validate URL format
	// For simplicity, let's just check if it's not empty
	if _, err := url.ParseRequestURI(prod.ImageURL); err != nil {
		return errors.New("image URL must be a valid URL")
	}

	// Notes validation
	if prod.Notes == "" || len(prod.Notes) > 200 {
		return errors.New("notes must not be empty and should be between 1 and 200 characters long")
	}

	// Price validation
	if prod.Price < 1 {
		return errors.New("price must be greater than or equal to 1")
	}

	// Stock validation
	if prod.Stock < 0 || prod.Stock > 100000 {
		return errors.New("stock must be between 0 and 100,000")
	}

	// Location validation
	if prod.Location == "" || len(prod.Location) > 200 {
		return errors.New("location must not be empty and should be between 1 and 200 characters long")
	}

	// isAvailable validation
	// Assuming isAvailable must be set explicitly
	if prod.IsAvailable == nil {
		return errors.New("isAvailable must be true or false")
	}

	// If all validations pass, return nil (no error)
	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}
