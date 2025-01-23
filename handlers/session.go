package handlers

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
)

var (
	store       *sessions.CookieStore
	sessionName string
)

func InitSession() {
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET is not set in .env file")
	}

	sessionName = os.Getenv("SESSION_NAME")
	if sessionName == "" {
		log.Fatal("SESSION_NAME is not set in .env file")
	}

	store = sessions.NewCookieStore([]byte(sessionSecret))
}

func GetStore() *sessions.CookieStore {
	return store
}

func GetSessionName() string {
	return sessionName
}
