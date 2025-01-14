package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Регистрация драйвера PostgreSQL
)

func InitDB() (*sqlx.DB, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}
