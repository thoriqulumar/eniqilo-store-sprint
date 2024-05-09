package repo

import (
	"context"
	"eniqilo-store/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CheckoutRepo interface {
	CreateCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error)
	GetCustomerById(ctx context.Context, userId string) (customer model.Customer, err error)
	GetCustomerByNumber(ctx context.Context, phoneNumber string) (customer model.Customer, err error)
	GetProductById(ctx context.Context, productId string) (product model.Product, err error)
	UpdateStockProduct(ctx context.Context, currentStock int, productId string) (err error)
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
	VALUES ($1, $2, $3, NOW())
	RETURNING *;`
)

func (r *checkoutRepo) CreateCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error) {

	customerId := uuid.New()

	err = r.db.QueryRowxContext(ctx, createCustomerQuery, customerId, data.PhoneNumber, data.Name).StructScan(&customer)

	return customer, err
}

var (
	getCustomerQuery = `SELECT * FROM "customer" WHERE "userId" = $1 LIMIT 1;`
)

func (r *checkoutRepo) GetCustomerById(ctx context.Context, userId string) (customer model.Customer, err error) {

	err = r.db.QueryRowxContext(ctx, getCustomerQuery, userId).StructScan(&customer)

	return customer, err
}

var (
	getCustomerByNumberQuery = `SELECT * FROM "customer" WHERE "phoneNumber" = $1 LIMIT 1;`
)

func (r *checkoutRepo) GetCustomerByNumber(ctx context.Context, phoneNumber string) (customer model.Customer, err error) {

	err = r.db.QueryRowxContext(ctx, getCustomerByNumberQuery, phoneNumber).StructScan(&customer)

	return customer, err
}

var (
	getProductQuery = `SELECT * FROM "product" WHERE "id" = $1 LIMIT 1;`
)

func (r *checkoutRepo) GetProductById(ctx context.Context, productId string) (product model.Product, err error) {

	err = r.db.QueryRowxContext(ctx, getProductQuery, productId).StructScan(&product)

	return product, err
}


var (
	updateStockProductQuery = `UPDATE FROM "product" SET "stock" = $1 WHERE id = $2;`
)

func (r *checkoutRepo) UpdateStockProduct(ctx context.Context, currentStock int, productId string) (err error) {

	_, err = r.db.ExecContext(ctx, updateStockProductQuery, productId)
	if err != nil {
		return err
	}
	return  nil
}

