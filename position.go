package main

import (
	"fmt"
	"strings"
)

type Position struct {
	turn  string
	board string
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
