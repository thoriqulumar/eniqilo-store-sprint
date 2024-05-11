package repo

import (
	"context"
	"eniqilo-store/model"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	GetProductByID(ctx context.Context, id uuid.UUID) (model.Product, error)
	GetProduct(ctx context.Context, param model.GetProductParam) ([]model.Product, error)
	CreateProduct(ctx context.Context, data model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, data model.Product) error
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

func (r *productRepo) CreateProduct(ctx context.Context, data model.Product) (model.Product, error) {

	// Generate UUID for the product ID
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return model.Product{}, fmt.Errorf("error generating UUID: %v", err)
	}
	data.ID = newUUID
	createdAt := time.Now()

	err = r.db.QueryRowxContext(ctx, createProductQuery,
		data.ID, data.Name, data.SKU, data.Category, data.ImageURL, data.Notes, data.Stock, data.Price, data.IsAvailable, data.Location, createdAt).Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		return model.Product{}, fmt.Errorf("error executing query: %v", err)
	}

	return data, nil
}

var updateProductQuery = `UPDATE product
SET "name"=$1, sku=$2, "category"=$3, stock=$4, price=$5, "imageUrl"=$6, notes=$7, "isAvailable"=$8, "location"=$9
WHERE id=$10;
`

func (r *productRepo) UpdateProduct(ctx context.Context, data model.Product) error {
	err := r.db.QueryRowxContext(ctx, updateProductQuery,
		data.Name, data.SKU, data.Category, data.Stock, data.Price, data.ImageURL, data.Notes, data.IsAvailable, data.Location, data.ID).Err()
	if err != nil {
		return err
	}

	return nil
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

func (r *productRepo) GetProductByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var product model.Product
	query := `SELECT * FROM product WHERE id = $1`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *productRepo) GetProduct(ctx context.Context, param model.GetProductParam) (products []model.Product, err error) {
	query := `SELECT * FROM product WHERE true ` + generateGetProductSQLFilter(param)
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return products, err
	}
	defer rows.Close()
	for rows.Next() {
		var product model.Product
		if err := rows.StructScan(&product); err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, nil
}

func generateGetProductSQLFilter(params model.GetProductParam) string {
	var conditions []string

	// Add conditions based on the fields provided
	if params.ID != nil {
		conditions = append(conditions, fmt.Sprintf(`"id" = '%s'`, *params.ID))
	}

	// TODO: explore searchable index
	if params.Name != nil {
		name := *params.Name
		// Append wildcard symbols to allow partial matching
		name = "%" + strings.ToLower(name) + "%"
		conditions = append(conditions, fmt.Sprintf(`lower("name") LIKE '%s'`, name))
	}
	if params.SKU != nil {
		conditions = append(conditions, fmt.Sprintf(`"sku" = '%s'`, *params.SKU))
	}
	if params.IsAvailable != nil {
		conditions = append(conditions, fmt.Sprintf(`"isAvailable" = %t`, *params.IsAvailable))
	}
	if params.Category != nil {
		conditions = append(conditions, fmt.Sprintf(`"category" = '%s'`, *params.Category))
	}
	if params.InStock != nil {
		conditions = append(conditions, fmt.Sprintf(`"stock" > 0`))
	}

	// Combine conditions with AND
	filter := strings.Join(conditions, " AND ")
	if filter != "" {
		// add and clause in the front
		filter = "AND " + filter
	}

	orderByClause := ""
	if params.Sort.Price != nil {
		orderByClause = fmt.Sprintf(" ORDER BY price %s", *params.Sort.Price)
	}

	// set default sort
	if params.Sort.CreatedAt == nil {
		defaultSort := "desc"
		params.Sort.CreatedAt = &defaultSort
	}
	if orderByClause != "" {
		orderByClause += fmt.Sprintf(`, "createdAt" %s`, *params.Sort.CreatedAt)
	} else {
		orderByClause = fmt.Sprintf(` ORDER BY "createdAt" %s`, *params.Sort.CreatedAt)
	}

	if orderByClause != "" {
		filter += " " + orderByClause
	}

	// Add additional clauses such as LIMIT and OFFSET
	if params.Limit != nil {
		filter += fmt.Sprintf(" LIMIT %d", *params.Limit)
	} else {
		filter += " LIMIT 5"
	}
	if params.Offset != nil {
		filter += fmt.Sprintf(" OFFSET %d", *params.Offset)
	} else {
		filter += " OFFSET 0"
	}

	return filter
}
