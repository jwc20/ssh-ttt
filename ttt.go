package ttt

import (
	"errors"
	"fmt"
	"strings"
)

type TicTacToe struct {
	store    PlayerStore
	position *Position
}

func NewTicTacToe(store PlayerStore) *TicTacToe {
	return &TicTacToe{
		store: store,
	}
}

func (g *TicTacToe) Start(numberOfPlayers int) {
	g.position = InitPosition()
}

func (g *TicTacToe) Finish(winner string) {
	g.store.RecordWin(winner)
}

func (g *TicTacToe) MakeMove(position int) error {
	idx := position - 1
	if g.position.Board[idx] != ' ' {
		return errors.New("square already taken")
	}
	g.position.Move(idx)
	return nil
}

func (g *TicTacToe) Board() string {
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		row := g.position.Board[i*3 : i*3+3]
		sb.WriteString(fmt.Sprintf(" %c | %c | %c \n", row[0], row[1], row[2]))
		if i < 2 {
			sb.WriteString("-----------\n")
		}
	}
	return sb.String()
}

func (g *TicTacToe) CurrentPlayer() string {
	return strings.ToUpper(g.position.Turn)
}

func (g *TicTacToe) IsOver() bool {
	return g.position.IsGameEnd()
}

func (g *TicTacToe) Winner() string {
	if g.position.IsWinFor("x") {
		return "X"
	}
	if g.position.IsWinFor("o") {
		return "O"
	}
	return ""
}

func (g *TicTacToe) BestMove() int {
	return g.position.BestMove() + 1
}
