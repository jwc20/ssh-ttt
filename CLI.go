package ttt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const PlayerPrompt = "Player %s, enter your move (1-9): "
const BadMoveInputErrMsg = "Bad value received for move, please enter a number between 1 and 9\n"
const SquareTakenErrMsg = "That square is already taken, please choose another\n"
const DrawMsg = "It's a draw!\n"
const WinMsg = "Player %s wins!\n"
const BoardHeader = "\nCurrent board:\n"

type TicTacToeGame interface {
	Game
	MakeMove(position int) error
	Board() string
	CurrentPlayer() string
	IsOver() bool
	Winner() string
}

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game TicTacToeGame
}

func NewCLI(in io.Reader, out io.Writer, game TicTacToeGame) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayGame() {
	cli.game.Start(2)

	for !cli.game.IsOver() {
		fmt.Fprint(cli.out, BoardHeader)
		fmt.Fprint(cli.out, cli.game.Board())
		fmt.Fprintf(cli.out, PlayerPrompt, cli.game.CurrentPlayer())

		input := cli.readLine()
		position, err := strconv.Atoi(input)

		if err != nil || position < 1 || position > 9 {
			fmt.Fprint(cli.out, BadMoveInputErrMsg)
			continue
		}

		err = cli.game.MakeMove(position)
		if err != nil {
			fmt.Fprint(cli.out, SquareTakenErrMsg)
			continue
		}
	}

	fmt.Fprint(cli.out, BoardHeader)
	fmt.Fprint(cli.out, cli.game.Board())

	winner := cli.game.Winner()
	if winner != "" {
		fmt.Fprintf(cli.out, WinMsg, winner)
		cli.game.Finish(winner)
	} else {
		fmt.Fprint(cli.out, DrawMsg)
	}
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
