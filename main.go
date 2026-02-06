package main

import (
	"fmt"
	"os"
)

func main() {
	game := initGame()

	fmt.Println("Tic-Tac-Toe")
	fmt.Println("Who goes first?")
	fmt.Println("1. Computer")
	fmt.Println("2. Human")
	fmt.Println("q. Quit")

	currentPlayer := game.askForPlayer(os.Stdin)
	if currentPlayer == -1 {
		fmt.Println("Goodbye!")
		return
	}

	for {
		// game.position.display()

		var move int
		if currentPlayer == COMPUTER {
			fmt.Println("Computer's turn")
			move = game.position.BestMove()
		} else {
			fmt.Println("Your turn (enter 1-9):")
			move = game.askForMove()
			if move == -1 {
				fmt.Println("Goodbye!")
				return
			}
		}

		game.position.Move(move)

		// if game.position.isWin() {
		// 	// game.position.display()

		// 	if currentPlayer == COMPUTER {
		// 		fmt.Println("Computer wins!")
		// 	} else {
		// 		fmt.Println("You win!")
		// 	}
		// 	break
		// }

		// if game.position.isDraw() {
		// 	game.position.display()
		// 	fmt.Println("It's a draw!")
		// 	break
		// }

		currentPlayer = 1 - currentPlayer
	}
}
