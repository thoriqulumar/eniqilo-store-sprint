package database

import (
	"eniqilo-store/config"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase() (*sqlx.DB, error) {
	params := strings.ReplaceAll(config.GetString("DB_PARAMS"), `"`, "")

	connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		config.GetString("DB_USERNAME"),
		config.GetString("DB_PASSWORD"),
		config.GetString("DB_HOST"),
		config.GetString("DB_PORT"),
		config.GetString("DB_NAME"),
		params,
	)

	db, err := sqlx.Open("postgres", connectionURL)

	return db, err
}
