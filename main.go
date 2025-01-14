package main

import "NotesWebApp/database"

func main() {
	pgdb, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer pgdb.Close()
}
