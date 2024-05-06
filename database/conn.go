package database

import (
	"eniqilo-store/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase() (*sqlx.DB, error) {
	connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		config.GetString("DB_USERNAME"),
		config.GetString("DB_PASSWORD"),
		config.GetString("DB_HOST"),
		config.GetString("DB_PORT"),
		config.GetString("DB_NAME"),
		config.GetString("DB_PARAMS"),
	)

	db, err := sqlx.Open("postgres", connectionURL)

	return db, err
}
