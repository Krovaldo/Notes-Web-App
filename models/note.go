package models

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Note struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (n *Note) CreateNote(db *sqlx.DB) error {
	query := `INSERT INTO notes (title, content, user_id) 
VALUES (:title, :content, :user_id) RETURNING id, created_at, updated_at`
	return db.QueryRowx(query, n).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func (n *Note) UpdateNote(db *sqlx.DB) error {
	query := `UPDATE notes SET title=:title, content=:content, updated_at=:updated_at 
             WHERE id=:id`
	_, err := db.NamedExec(query, n)
	return err
}

func (n *Note) DeleteNote(db *sqlx.DB) error {
	query := `DELETE FROM notes WHERE id=:id`
	_, err := db.NamedExec(query, n)
	return err
}
func (n *Note) GetNotesByUser(db *sqlx.DB, userId int) ([]Note, error) {
	var notes []Note
	query := `SELECT id, title, content, created_at, updated_at FROM notes WHERE user_id=:user_id`
	err := db.Select(&notes, query, userId)
	return notes, err
}

func GetNoteByID(db *sqlx.DB, id int) (*Note, error) {
	var note Note
	query := `SELECT id, title, content, created_at, updated_at FROM notes WHERE id=$1`

	err := db.Get(&note, query, id)
	return &note, err
}
