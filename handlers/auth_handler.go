package handlers

import (
	"NotesWebApp/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type AuthHandler struct {
	DB *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (ah *AuthHandler) Index(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)
	}
	_, ok := session.Values["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusFound)
}

func (ah *AuthHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error while executing login.html template:", err)
		return
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(ah.DB, email)
	if err != nil {
		http.Error(w, "User not found: invalid credentials", http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Printf("Failed to get session: %v", err)

		http.Error(w, "Failed to get session", http.StatusBadRequest)
		return
	}
	session.Values["userID"] = user.ID
	err = session.Save(r, w)
	if err != nil {
		log.Println("Can't save session:", err)
		return
	}
	log.Println("User logged in successfully:", user.ID) //123
	http.Redirect(w, r, "/notes", http.StatusFound)
}

func (ah *AuthHandler) RegisterForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error while executing register.html template:", err)
		return
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password for user %s: %v", email, err)

		http.Error(w, "Internal server error: failed to process password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := user.CreateUser(ah.DB); err != nil {
		log.Printf("Failed to create user %s: %v", email, err)

		http.Error(w, "Internal server error: failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("user created %+v", user)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	session.Values = make(map[interface{}]interface{}) // удаляем все
	session.Options.MaxAge = -1                        // ставим срок действия сессии в прошлое (удаление сессии)

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
