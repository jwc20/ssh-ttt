package main

import (
	"fmt"
	"strings"
)

type Position struct {
	turn  string
	board string
}

func (p Position) Move(i int) Position {
	p.board = p.board[:i] + p.turn + p.board[i+1:]

	switch p.turn {
	case "x":
		p.turn = "o"
	case "o":
		p.turn = "x"
	}
	return p
}

func (p Position) String() string {
	return fmt.Sprintf("%s.%s", p.turn, p.board)
}

func initPosition() Position {
	return Position{
		turn:  "x",
		board: strings.Repeat(" ", 9),
	}
}
