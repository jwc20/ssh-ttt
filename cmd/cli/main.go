package main

import (
	"fmt"
	"log"
	"os"

	ttt "github.com/jwc20/ssh-ttt"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := ttt.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play tic-tac-toe")
	fmt.Println("Type {Name} wins to record a win")

	ttt.NewCLI(store, os.Stdin).PlayTTT()
}
