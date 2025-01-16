package models

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *User) CreateUser(db *sqlx.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT into users (email, password) VALUES ($1, $2) RETURNING id`
	return db.QueryRowx(query, u.Email, string(hashedPassword)).Scan(&u.ID)
}

func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
	var user User
	query := `SELECT id, email, password FROM users WHERE email=$1`
	err := db.Get(&user, query, email)
	return &user, err
}
