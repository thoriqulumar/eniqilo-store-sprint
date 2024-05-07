package repo

import (
	"context"
	"eniqilo-store/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CheckoutRepo interface {
	CreateCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error)
}

type checkoutRepo struct {
	db *sqlx.DB
}

func NewCheckoutRepo(db *sqlx.DB) CheckoutRepo {
	return &checkoutRepo{
		db: db,
	}
}

var (
	createCustomerQuery = `INSERT INTO "customer" ("userId", "phoneNumber", "name", "createdAt") 
	VALUES ($1, $2, $3, NOW())`
)

func (r *checkoutRepo) CreateCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error) {

	customerId := uuid.New()

	err = r.db.QueryRowxContext(ctx, createCustomerQuery, customerId, data.PhoneNumber, data.Name).StructScan(&customer)

	return customer, err
}
