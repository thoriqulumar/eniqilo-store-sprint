package repo

import (
	"context"
	"eniqilo-store/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	CreateProduct(ctx context.Context, data model.CreatedProduct) (model.CreatedProduct, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) ProductRepo {
	return &productRepo{db: db}
}

var createProductQuery = `INSERT INTO product 
    ("id",name, sku, category, "imageUrl", notes, stock, price, "isAvailable", location, "createdAt")
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    RETURNING "id", "createdAt"`

func (r *productRepo) CreateProduct(ctx context.Context, data model.CreatedProduct) (model.CreatedProduct, error) {

	// Generate UUID for the product ID
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return model.CreatedProduct{}, fmt.Errorf("error generating UUID: %v", err)
	}
	data.ID = newUUID
	createdAt := time.Now()

	err = r.db.QueryRowxContext(ctx, createProductQuery,
		data.ID, data.Name, data.SKU, data.Category, data.ImageURL, data.Notes, data.Stock, data.Price, data.IsAvailable, data.Location, createdAt).Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		return model.CreatedProduct{}, fmt.Errorf("error executing query: %v", err)
	}

	return data, nil
}

func (r *productRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM product WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}
