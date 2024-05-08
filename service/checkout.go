package service

import (
	"context"
	"database/sql"
	"eniqilo-store/model"
	"eniqilo-store/repo"
	cerr "eniqilo-store/utils/error"
	"errors"
	"net/http"
)

type CheckoutService interface {
	CreateNewCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error)
	ValidateUser(ctx context.Context, userId string) (customer model.Customer, err error) 
	ValidateProduct(ctx context.Context, products []model.ProductDetail) ( total float32, err error)
	CheckoutProduct(ctx context.Context, products []model.ProductDetail) (err error)
}

type checkoutService struct {
	repo repo.CheckoutRepo
}

func NewCheckoutService(r repo.CheckoutRepo) CheckoutService {
	return &checkoutService{
		repo: r,
	}
}

func (s *checkoutService) CreateNewCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error) {

	dataCustomer, err := s.repo.CreateCustomer(ctx, data)
	if err != nil {
		return
	}

	return dataCustomer, nil
}

func (s *checkoutService) ValidateUser(ctx context.Context, userId string) (customer model.Customer, err error) {

	dataCustomer, err := s.repo.GetCustomerById(ctx, userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model.Customer{}, cerr.New(http.StatusNotFound, "productId is not found")
	}

	return dataCustomer, nil
}

func (s *checkoutService) ValidateProduct(ctx context.Context, products []model.ProductDetail) ( total float32, err error) {
	var totalPrice float32
	for _, product := range products {
		dataProduct, err := s.repo.GetProductById(ctx, product.ProductId)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return 0, cerr.New(http.StatusNotFound, "productId is not found")
		}

		if dataProduct.Stock < product.Quantity{
			return 0, cerr.New(http.StatusBadRequest, `quantity product id`+product.ProductId+` is not enough`)
		}

		if !dataProduct.IsAvailable{
			return 0, cerr.New(http.StatusBadRequest, `quantity product id`+product.ProductId+` is not available`)
		}

		totalPrice += (float32(dataProduct.Price) * float32(product.Quantity) )
	}

	return totalPrice, nil
}

func (s *checkoutService) CheckoutProduct(ctx context.Context, products []model.ProductDetail) (err error){
	for _, product := range products {
		dataProduct, err := s.repo.GetProductById(ctx, product.ProductId)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return cerr.New(http.StatusNotFound, "productId is not found")
		}

		currentStock := dataProduct.Stock - product.Quantity
		err = s.repo.UpdateStockProduct(ctx, currentStock, product.ProductId)
		if err != nil{
			return cerr.New(http.StatusInternalServerError, "error when update product" + product.ProductId)
		}
		
	}

	return nil
}