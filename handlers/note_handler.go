package handlers

import (
	"NotesWebApp/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type NoteHandler struct {
	DB *sqlx.DB
}

func NewNoteHandler(db *sqlx.DB) *NoteHandler {
	return &NoteHandler{DB: db}
}

func (nh *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	note := models.Note{}

	notes, err := note.GetNotesByUser(nh.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err = tmpl.Execute(w, notes)
	if err != nil {
		log.Println("Error while executing index.html:", err)
		return
	}
}

func (nh *NoteHandler) CreateNoteForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error while executing create.html:", err)
		return
	}
}

func (nh *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)
		return
	}
	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	note := &models.Note{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := note.CreateNote(nh.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusFound)
}

func (nh *NoteHandler) EditNoteForm(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	note, err := models.GetNoteByID(nh.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("полученная заметка GetNoteByID: %+v", note)

	if note.UserID != userID {
		log.Printf("User %d tried to edit note %d belonging to user %d", userID, note.ID, note.UserID)
		http.Error(w, "You do not have permission to edit this note", http.StatusForbidden)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	err = tmpl.Execute(w, note)
	if err != nil {
		log.Println("Error while executing edit.html:", err)
		return
	}
}

func (nh *NoteHandler) EditNote(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := models.GetNoteByID(nh.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else {
		if note.UserID != userID {
			log.Printf("User %d tried to edit note %d belonging to user %d", userID, note.ID, note.UserID)
			http.Error(w, "You do not have permission to edit this note", http.StatusForbidden)
			return
		}
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	note.Title = title
	note.Content = content

	if err := note.UpdateNote(nh.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusFound)
}

func (nh *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := models.GetNoteByID(nh.DB, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	} else {
		if note.UserID != userID {
			log.Printf("User %d tried to edit note %d belonging to user %d", userID, note.ID, note.UserID)
			http.Error(w, "You do not have permission to edit this note", http.StatusForbidden)
			return
		}
	}

	if err := note.DeleteNote(nh.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}
