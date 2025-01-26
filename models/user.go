package models

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *User) CreateUser(db *sqlx.DB) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	return db.QueryRowx(query, u.Email, u.Password).Scan(&u.ID)
}

func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
	var user User
	query := `SELECT id, email, password FROM users WHERE email=$1`
	err := db.Get(&user, query, email)
	return &user, err
}
