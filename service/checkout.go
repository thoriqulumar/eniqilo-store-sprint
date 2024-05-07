package service

import (
	"context"
	"eniqilo-store/model"
	"eniqilo-store/repo"
)

type CheckoutService interface {
	CreateNewCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error)
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
