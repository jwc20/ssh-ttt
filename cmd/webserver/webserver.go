package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jwc20/ssh-ttt/handlers"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		public_key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	createRoomsTable := `CREATE TABLE IF NOT EXISTS rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		winner TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		finished_at DATETIME
	);`

	_, err = db.Exec(createRoomsTable)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	router := gin.Default()

	router.GET("/users", handlers.ListUser(db))
	router.POST("/users", handlers.CreateUser(db))

	port := ":8080"
	log.Printf("Server is running on http://localhost%s", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
