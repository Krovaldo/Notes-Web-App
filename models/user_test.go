package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestUser_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	user := &User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`)).
		WithArgs(user.Email, user.Password).
		WillReturnRows(rows)

	err = user.CreateUser(sqlxDB)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	email := "test@example.com"
	expectedUser := &User{
		ID:       1,
		Email:    email,
		Password: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password)

	mock.ExpectQuery(`SELECT id, email, password FROM users WHERE email=\$1`).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := GetUserByEmail(sqlxDB, email)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}
