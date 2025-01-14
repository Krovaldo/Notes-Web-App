package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func InitDB() (*sqlx.DB, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}
