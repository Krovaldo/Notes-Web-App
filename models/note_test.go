package models

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNote_CreateNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	note := &Note{
		Title:   "Test Title",
		Content: "Test Content",
		UserID:  1,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO notes (title, content, user_id) 
VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`)).
		WithArgs(note.Title, note.Content, note.UserID).
		WillReturnRows(rows)

	err = note.CreateNote(sqlxDB)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNote_UpdateNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Faild to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	note := &Note{
		ID:        1,
		Title:     "Updated Title",
		Content:   "Updated Content",
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE notes SET title=?, content=?, updated_at=?
             WHERE id=?`)).
		WithArgs(note.Title, note.Content, sqlmock.AnyArg(), note.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = note.UpdateNote(sqlxDB)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNote_DeleteNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	note := &Note{
		ID: 1,
	}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM notes WHERE id=?`)).WithArgs(note.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = note.DeleteNote(sqlxDB)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNotesByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	userID := 1
	expectedNotes := []Note{
		{ID: 1, Title: "Note 1", Content: "Content 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Note 2", Content: "Content 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(expectedNotes[0].ID, expectedNotes[0].Title, expectedNotes[0].Content, expectedNotes[0].CreatedAt,
			expectedNotes[0].UpdatedAt).
		AddRow(expectedNotes[1].ID, expectedNotes[1].Title, expectedNotes[1].Content, expectedNotes[1].CreatedAt,
			expectedNotes[1].UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, content, created_at, updated_at 
FROM notes WHERE user_id=$1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	note := &Note{}

	notes, err := note.GetNotesByUser(sqlxDB, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedNotes, notes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNoteByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	noteID := 1
	expectedNote := &Note{
		ID:        noteID,
		Title:     "Test Note",
		Content:   "Test Content",
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"}).
		AddRow(expectedNote.ID, expectedNote.Title, expectedNote.Content, expectedNote.UserID,
			time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at 
FROM notes WHERE id=$1`)).
		WithArgs(noteID).
		WillReturnRows(rows)

	note, err := GetNoteByID(sqlxDB, noteID)
	assert.NoError(t, err)
	assert.NotNil(t, note)

	assert.Equal(t, expectedNote.ID, noteID)
	assert.Equal(t, expectedNote.Title, note.Title)
	assert.Equal(t, expectedNote.Content, note.Content)
	assert.Equal(t, expectedNote.UserID, note.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNoteByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	noteID := 1
	expectedError := errors.New("database connection error")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at 
FROM notes WHERE id=$1`)).
		WithArgs(noteID).
		WillReturnError(expectedError)

	note, err := GetNoteByID(sqlxDB, noteID)
	assert.Error(t, err)
	assert.Nil(t, note)
	assert.Equal(t, expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNoteByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	userID := 1
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at 
FROM notes WHERE id=$1`)).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	note, err := GetNoteByID(sqlxDB, userID)
	assert.NoError(t, err)
	assert.Nil(t, note)

	assert.NoError(t, mock.ExpectationsWereMet())
}
