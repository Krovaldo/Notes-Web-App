package main

import (
	"NotesWebApp/database"
	"NotesWebApp/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

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
	defer db.Close()

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

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
