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

//
//import (
//	"fmt"
//	"os"
//)
//
//func main() {
//	game := initGame()
//
//	fmt.Println("Tic-Tac-Toe")
//	fmt.Println("Who goes first?")
//	fmt.Println("1. Computer")
//	fmt.Println("2. Human")
//	fmt.Println("q. Quit")
//
//	currentPlayer := game.askForPlayer(os.Stdin)
//	if currentPlayer == -1 {
//		fmt.Println("Goodbye!")
//		return
//	}
//
//	for {
//		// game.position.display()
//
//		var move int
//		if currentPlayer == COMPUTER {
//			fmt.Println("Computer's turn")
//			move = game.position.BestMove()
//		} else {
//			fmt.Println("Your turn (enter 1-9):")
//			move = game.askForMove()
//			if move == -1 {
//				fmt.Println("Goodbye!")
//				return
//			}
//		}
//
//		game.position.Move(move)
//
//		// if game.position.isWin() {
//		// 	// game.position.display()
//
//		// 	if currentPlayer == COMPUTER {
//		// 		fmt.Println("Computer wins!")
//		// 	} else {
//		// 		fmt.Println("You win!")
//		// 	}
//		// 	break
//		// }
//
//		// if game.position.isDraw() {
//		// 	game.position.display()
//		// 	fmt.Println("It's a draw!")
//		// 	break
//		// }
//
//		currentPlayer = 1 - currentPlayer
//	}
//}
