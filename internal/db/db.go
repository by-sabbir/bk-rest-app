package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataBase struct {
	Client *sqlx.DB
}

func NewDatabase() (*DataBase, error) {
	connecttionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	dbConn, err := sqlx.Connect("postgres", connecttionString)

	if err != nil {
		return &DataBase{}, fmt.Errorf("db init failed: %w", err)
	}
	return &DataBase{
		Client: dbConn,
	}, nil
}
