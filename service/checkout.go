package service

import (
	"context"
	"database/sql"
	"eniqilo-store/model"
	"eniqilo-store/repo"
	cerr "eniqilo-store/utils/error"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type CheckoutService interface {
	CreateNewCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error)
	ValidateUser(ctx context.Context, userId string) (customer model.Customer, err error)
	ValidateProduct(ctx context.Context, products []model.ProductDetail) (total float32, err error)
	CheckoutProduct(ctx context.Context, transaction model.Transaction) (err error)
	GetAllCustomer(ctx context.Context, name, phoneNumber string, limit, offset int) (listCustomer []model.CustomerResponseData, err error)
	GetAllTransaction(ctx context.Context, params model.GetHistoryParam) (listTransaction []model.Transaction, err error)
}

type checkoutService struct {
	repo   repo.CheckoutRepo
	logger *zap.Logger
}

func NewCheckoutService(r repo.CheckoutRepo, logger *zap.Logger) CheckoutService {
	return &checkoutService{
		repo:   r,
		logger: logger,
	}
}

func (s *checkoutService) CreateNewCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error) {
	if data.PhoneNumber == nil {
		return model.Customer{}, cerr.New(http.StatusBadRequest, "phoneNumber is required")
	}
	_, err = s.repo.GetCustomerByNumber(ctx, *data.PhoneNumber)
	if !errors.Is(err, sql.ErrNoRows) {
		return model.Customer{}, cerr.New(http.StatusConflict, "phoneNumber already exists")
	}

	dataCustomer, err := s.repo.CreateCustomer(ctx, data)
	if err != nil {
		return model.Customer{}, cerr.New(http.StatusInternalServerError, "Internal Server Error")
	}

	return dataCustomer, nil
}

func (s *checkoutService) ValidateUser(ctx context.Context, userId string) (customer model.Customer, err error) {

	dataCustomer, err := s.repo.GetCustomerById(ctx, userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model.Customer{}, cerr.New(http.StatusNotFound, "customerId is not found")
	}

	if dataCustomer.UserId == "" {
		return model.Customer{}, cerr.New(http.StatusNotFound, "customerId is not found")
	}
	return dataCustomer, nil
}

func (s *checkoutService) ValidateProduct(ctx context.Context, products []model.ProductDetail) (total float32, err error) {
	var totalPrice float32
	for _, product := range products {
		if product.ProductId == "" {
			return totalPrice, cerr.New(http.StatusBadRequest, "productId is required")
		}

		dataProduct, err := s.repo.GetProductById(ctx, product.ProductId)
		if err != nil {
			return 0, cerr.New(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, sql.ErrNoRows) {
			return 0, cerr.New(http.StatusNotFound, "productId is not found")
		}

		if *dataProduct.Stock < product.Quantity {
			return 0, cerr.New(http.StatusBadRequest, `quantity product id `+product.ProductId+` is not enough`)
		}

		if !*dataProduct.IsAvailable {
			return 0, cerr.New(http.StatusBadRequest, `quantity product id `+product.ProductId+` is not available`)
		}

		totalPrice += (float32(dataProduct.Price) * float32(product.Quantity))
	}

	return totalPrice, nil
}

func (s *checkoutService) CheckoutProduct(ctx context.Context, transaction model.Transaction) (err error) {
	productIDs := make([]string, 0, len(transaction.ProductDetails))
	for _, product := range transaction.ProductDetails {
		if product.ProductId == "" {
			return cerr.New(http.StatusBadRequest, "productId cannot be empty")
		}
		productIDs = append(productIDs, product.ProductId)
	}

	productStocks, err := s.repo.GetProductStocks(ctx, productIDs)
	if err != nil {
		return cerr.New(http.StatusInternalServerError, "error fetching product stock levels: "+err.Error())
	}

	updatedStocks := make(map[string]int)
	for _, product := range transaction.ProductDetails {
		existingStock, ok := productStocks[product.ProductId]
		if !ok {
			// Handle unexpected missing product (shouldn't occur after previous check)
			return cerr.New(http.StatusInternalServerError, fmt.Sprintf("unexpected error: product %s not found in fetched stocks", product.ProductId))
		}
		updatedStocks[product.ProductId] = existingStock - product.Quantity
	}

	for productId, stock := range updatedStocks {
		err = s.repo.UpdateStockProduct(ctx, stock, productId)
		if err != nil {
			return cerr.New(http.StatusInternalServerError, fmt.Sprintf("error updating stock for product %s: %s", productId, err.Error()))
		}
	}

	err = s.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		return cerr.New(http.StatusInternalServerError, fmt.Sprintf("error inserting transaction data"))
	}

	return nil
}

func (s *checkoutService) GetAllCustomer(ctx context.Context, name, phoneNumber string, limit, offset int) (listCustomer []model.CustomerResponseData, err error) {
	dataCustomer, err := s.repo.GetAllCustomer(ctx, name, phoneNumber, limit, offset)
	if err != nil {
		s.logger.Error("failed getAllCustomer", zap.Error(err))
		return []model.CustomerResponseData{}, cerr.New(http.StatusInternalServerError, "Internal Server Error")
	}

	return dataCustomer, nil
}

func (s *checkoutService) GetAllTransaction(ctx context.Context, params model.GetHistoryParam) (listTransaction []model.Transaction, err error) {
	listTransaction, err = s.repo.GetHistoryTransaction(ctx, params)
	if err != nil {
		return
	}

	return listTransaction, nil
}
