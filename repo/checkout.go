package repo

import (
	"context"
	"encoding/json"
	"eniqilo-store/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CheckoutRepo interface {
	NewTx() (*sqlx.Tx, error)
	CreateCustomer(ctx context.Context, data model.CustomerRequest) (customer model.Customer, err error)
	GetCustomerById(ctx context.Context, userId string) (customer model.Customer, err error)
	GetCustomerByNumber(ctx context.Context, phoneNumber string) (customer model.Customer, err error)
	GetProductById(ctx context.Context, productId string) (product model.Product, err error)
	UpdateStockProduct(ctx context.Context, tx *sqlx.Tx, currentStock int, productId string) (err error)
	GetProductStocks(ctx context.Context, productIDs []string) (map[string]int, error)
	GetAllCustomer(ctx context.Context, name, phoneNumber string, limit, offset int) (customers []model.CustomerResponseData, err error)
	CreateTransaction(ctx context.Context, tx *sqlx.Tx, transaction model.Transaction) (err error)
	GetHistoryTransaction(ctx context.Context, params model.GetHistoryParam) (customers []model.Transaction, err error)
}

type checkoutRepo struct {
	db *sqlx.DB
}

func NewCheckoutRepo(db *sqlx.DB) CheckoutRepo {
	return &checkoutRepo{
		db: db,
	}
}

func (r *checkoutRepo) NewTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
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
	if err != nil {
		return
	}
	return product, nil
}

var (
	updateStockProductQuery = `UPDATE "product" SET "stock" = $1 WHERE id = $2;`
)

func (r *checkoutRepo) UpdateStockProduct(ctx context.Context, tx *sqlx.Tx, currentStock int, productId string) (err error) {
	_, err = tx.ExecContext(ctx, updateStockProductQuery, currentStock, productId)
	if err != nil {
		return err
	}
	return nil
}

var (
	getProductStocksQuery = `SELECT id, stock FROM product WHERE id = ANY ($1);`
)

func (r *checkoutRepo) GetProductStocks(ctx context.Context, productIDs []string) (map[string]int, error) {
	// Construct efficient query to fetch product IDs and stocks in one go
	rows, err := r.db.QueryContext(ctx, getProductStocksQuery, pq.Array(productIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productStocks := make(map[string]int)
	for rows.Next() {
		var id string
		var stock int
		if err := rows.Scan(&id, &stock); err != nil {
			return nil, err
		}
		productStocks[id] = stock
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productStocks, nil
}

func (r *checkoutRepo) GetAllCustomer(ctx context.Context, name, phoneNumber string, limit, offset int) (customers []model.CustomerResponseData, err error) {
	var listCustomer []model.CustomerResponseData
	var getAllCustomerQuery = `SELECT "userId", "phoneNumber", "name" FROM customer WHERE 1=1`

	if phoneNumber != "" {
		getAllCustomerQuery += fmt.Sprintf(` AND "phoneNumber" LIKE '%%%s%%'`, phoneNumber)
	}
	if name != "" {
		getAllCustomerQuery += fmt.Sprintf(` AND LOWER(name) LIKE LOWER('%%%s%%')`, name)
	}
	getAllCustomerQuery += fmt.Sprintf(` ORDER BY "createdAt" DESC LIMIT %d OFFSET %d`, limit, offset)

	rows, err := r.db.QueryContext(ctx, getAllCustomerQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var customer model.CustomerResponseData
		if err := rows.Scan(&customer.UserId, &customer.PhoneNumber, &customer.Name); err != nil {
			return nil, err
		}
		listCustomer = append(listCustomer, customer)
	}
	return listCustomer, nil
}

var (
	createTransactionQuery = `INSERT INTO "transaction" ("transactionId", "customerId", "productDetails", "paid", "change", "createdAt") VALUES ($1, $2, $3, $4, $5, NOW());`
)

func (r *checkoutRepo) CreateTransaction(ctx context.Context, tx *sqlx.Tx, transaction model.Transaction) (err error) {
	productDetailsByte, _ := json.Marshal(transaction.ProductDetails)
	_, err = tx.ExecContext(ctx, createTransactionQuery, transaction.TransactionId, transaction.CustomerId, productDetailsByte, transaction.Paid, transaction.Change)
	if err != nil {
		return err
	}

	return nil
}

func (r *checkoutRepo) GetHistoryTransaction(ctx context.Context, params model.GetHistoryParam) (customers []model.Transaction, err error) {
	var listTransaction []model.Transaction
	var getAllHistoryTransactionQuery = `SELECT * FROM "transaction" WHERE 1=1`

	if params.CustomerId != nil {
		getAllHistoryTransactionQuery += fmt.Sprintf(` AND "customerId" = %s`, params.CustomerId)
	}
	if params.CreatedAt != nil {
		if *params.CreatedAt != "desc" && *params.CreatedAt != "asc" {
			*params.CreatedAt = "desc"
		}
		getAllHistoryTransactionQuery += fmt.Sprintf(` ORDER BY "createdAt" %s`, *params.CreatedAt)
	} else {
		getAllHistoryTransactionQuery += ` ORDER BY "createdAt" DESC`
	}

	if params.Limit == 0 {
		params.Limit = 5 // default limit
	}

	getAllHistoryTransactionQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, getAllHistoryTransactionQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var transaction model.Transaction
		var productDetailsByte []byte
		if err := rows.Scan(&transaction.TransactionId, &transaction.CustomerId, &productDetailsByte, &transaction.Paid, &transaction.Change, &transaction.CreatedAt); err != nil {
			return nil, err
		}

		json.Unmarshal(productDetailsByte, &transaction.ProductDetails)

		listTransaction = append(listTransaction, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listTransaction, nil
}
