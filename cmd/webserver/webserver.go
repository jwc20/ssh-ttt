package main

import (
	"log"
	"net/http"
	"os"

	ttt "github.com/jwc20/ssh-ttt"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := ttt.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := ttt.NewPlayerServer(store)

	if err := http.ListenAndServe(":5002", server); err != nil {
		log.Fatalf("could not listen on port 5002 %v", err)
	}
}
