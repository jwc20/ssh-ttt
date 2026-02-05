package main

import (
	"fmt"
	"strings"
)

const (
	DIM  = 3
	SIZE = DIM * DIM
)

type Position struct {
	turn  string
	board string
}

func (p Position) choose(x, o string) string {
	if p.turn == "x" {
		return x
	}
	return o
}

func (p Position) Move(i int) Position {
	p.board = p.board[:i] + p.turn + p.board[i+1:]
	p.turn = p.choose("o", "x")
	return p
}

func (p Position) PossibleMoves() []int {
	var result []int
	for i, piece := range strings.Split(p.board, "") {
		if piece == " " {
			result = append(result, i)
		}
	}
	return result
}

func (p Position) isWinFor(piece string) bool {
	isMatch := func(line string) bool {
		return strings.Count(line, piece) == DIM
	}

	// 1. Check Rows
	for i := 0; i < SIZE; i += DIM {
		if isMatch(p.board[i : i+DIM]) {
			return true
		}
	}

	// 2. Check Columns
	for i := range DIM {
		col := ""
		for j := i; j < SIZE; j += DIM {
			col += string(p.board[j])
		}
		if isMatch(col) {
			return true
		}
	}

	// 3. Major Diagonal: Index sequence [0, DIM+1, 2(DIM+1), ...]
	majDiag := ""
	for i := 0; i < SIZE; i += DIM + 1 {
		majDiag += string(p.board[i])
	}
	if isMatch(majDiag) {
		return true
	}

	// 4. Minor Diagonal: Index sequence [DIM-1, 2(DIM-1), ..., SIZE-DIM]
	minDiag := ""
	for i := DIM - 1; i <= SIZE-DIM; i += DIM - 1 {
		minDiag += string(p.board[i])
	}
	if isMatch(minDiag) {
		return true
	}

	return false
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
