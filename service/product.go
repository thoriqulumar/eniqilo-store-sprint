package service

import (
	"context"
	"database/sql"
	"eniqilo-store/model"
	"eniqilo-store/repo"
	cerr "eniqilo-store/utils/error"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

// ProductService handles business logic related to products.
type ProductService interface {
	GetProduct(ctx context.Context, param model.GetProductParam) ([]model.Product, error)
	GetProductCustomer(ctx context.Context, param model.GetProductParam) ([]model.Product, error)
	CreateProduct(ctx context.Context, data model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, data model.Product) (model.Product, error)
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
func (s *productService) CreateProduct(ctx context.Context, prod model.Product) (model.Product, error) {
	// Validate the product
	if err := validateCreateProduct(prod); err != nil {
		return model.Product{}, cerr.New(http.StatusBadRequest, err.Error())
	}

	return s.repo.CreateProduct(ctx, prod)
}

// UpdateProduct handles the creation of a new product.
func (s *productService) UpdateProduct(ctx context.Context, prod model.Product) (model.Product, error) {
	// Validate the product
	if err := validateCreateProduct(prod); err != nil {
		return model.Product{}, cerr.New(http.StatusBadRequest, err.Error())
	}

	_, err := s.repo.GetProductByID(ctx, prod.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, cerr.New(http.StatusNotFound, "Product not found")
		}
		return model.Product{}, cerr.New(http.StatusInternalServerError, "Error updating product")
	}

	return prod, s.repo.UpdateProduct(ctx, prod)
}

func validateCreateProduct(prod model.Product) error {
	// Name validation
	if prod.Name == "" || len(prod.Name) > 30 {
		return errors.New("name must not be empty and should be between 1 and 30 characters long")
	}

	// SKU validation
	if prod.SKU == "" || len(prod.SKU) > 30 {
		return errors.New("SKU must not be empty and should be between 1 and 30 characters long")
	}

	// Category validation
	validCategories := map[model.Category]struct{}{
		model.Clothing:    {},
		model.Accessories: {},
		model.Footwear:    {},
		model.Beverages:   {},
	}
	if _, ok := validCategories[model.Category(prod.Category)]; !ok {
		return errors.New("invalid category")
	}

	// Image URL validation
	// You can use regex or a library like net/url to validate URL format
	// For simplicity, let's just check if it's not empty
	ok := isValidURL(prod.ImageURL)
	if !ok {
		return errors.New("url is invalid")
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
	if prod.Stock == nil {
		return errors.New("stock must not be nil")
	}
	if *prod.Stock < 0 || *prod.Stock > 100000 {
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

func (s *productService) GetProduct(ctx context.Context, param model.GetProductParam) (product []model.Product, err error) {
	// generate filter query from param
	// do request
	product, err = s.repo.GetProduct(ctx, param)
	if product == nil || err != nil {
		product = []model.Product{}
	}
	return product, err
}

func (s *productService) GetProductCustomer(ctx context.Context, param model.GetProductParam) ([]model.Product, error) {
	*param.IsAvailable = true
	return s.repo.GetProduct(ctx, param)
}
