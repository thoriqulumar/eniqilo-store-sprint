package repo

import (
	"eniqilo-store/model"

	"github.com/jmoiron/sqlx"
)

type Repo interface{}

type repo struct {
	db *sqlx.DB
}

type ProductRepository interface {
	Register(*model.Product) (*model.Product, error)
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}
