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

	fmt.Println("Let's play Tic-Tac-Toe!")

	game := ttt.NewTicTacToe(store)
	cli := ttt.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayGame()
}
