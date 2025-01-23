package main

import (
	"log"
	"net/http"
	"time"

	"NotesWebApp/database"
	"NotesWebApp/handlers"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func runServer(db *sqlx.DB, server *http.Server) error {
	defer db.Close()

	log.Println("Server is running on port 8080...")
	return server.ListenAndServe()
}

func main() {
	err := godotenv.Load(".env") // Загружаем переменные окружения
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handlers.InitSession() // Инициализация сессии

	db, err := database.InitDB() // инициализация базы
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter() // инициализация роутера

	// инициализация обработчиков
	noteHandler := handlers.NewNoteHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// маршруты заметок
	router.HandleFunc("/notes", noteHandler.GetNotes).Methods("GET")
	router.HandleFunc("/notes/create", noteHandler.CreateNoteForm).Methods("GET")
	router.HandleFunc("/notes/create", noteHandler.CreateNote).Methods("POST")
	router.HandleFunc("/notes/edit/{id}", noteHandler.EditNoteForm).Methods("GET")
	router.HandleFunc("/notes/edit/{id}", noteHandler.EditNote).Methods("POST")
	router.HandleFunc("/notes/delete/{id}", noteHandler.DeleteNote).Methods("POST")

	// маршруты аутентификации
	router.HandleFunc("/", authHandler.Index).Methods("GET")
	router.HandleFunc("/login", authHandler.LoginForm).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.RegisterForm).Methods("GET")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Server is running on port 8080...")
	if err := runServer(db, server); err != nil {
		log.Fatal("Server error:", err)
	}
}
