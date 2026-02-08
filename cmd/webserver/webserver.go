package main

import (
	"log"

	"github.com/gin-gonic/gin"
	ttt "github.com/jwc20/ssh-ttt"
	"github.com/jwc20/ssh-ttt/handlers"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "/.app.db"

func main() {
	store, close, err := ttt.FileSystemTTTStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	db := store.Database

	router := gin.Default()
	router.GET("/users", handlers.ListUser(db))
	router.POST("/users", handlers.CreateUser(db))

	port := ":8080"
	log.Printf("Server is running on http://localhost%s", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
