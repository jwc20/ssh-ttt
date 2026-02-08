package ttt

import (
	"fmt"
	"math"
	"strings"
)

const (
	DIM  = 3
	SIZE = DIM * DIM
)

type Position struct {
	Turn  string
	Board string
}

func InitPosition() *Position {
	return &Position{
		Turn:  "x",
		Board: strings.Repeat(" ", 9),
	}
}

func (p Position) choose(x, o string) string {
	if p.Turn == "x" {
		return x
	}
	return o
}

func (p *Position) Move(i int) *Position {
	p.Board = p.Board[:i] + p.Turn + p.Board[i+1:]
	p.Turn = p.choose("o", "x")
	return p
}

func (p Position) PossibleMoves() []int {
	var result []int
	for i, piece := range strings.Split(p.Board, "") {
		if piece == " " {
			result = append(result, i)
		}
	}
	return result
}

func (p Position) IsWinFor(piece string) bool {
	isMatch := func(line string) bool {
		return strings.Count(line, piece) == DIM
	}

	for i := 0; i < SIZE; i += DIM {
		if isMatch(p.Board[i : i+DIM]) {
			return true
		}
	}

	for i := range DIM {
		col := ""
		for j := i; j < SIZE; j += DIM {
			col += string(p.Board[j])
		}
		if isMatch(col) {
			return true
		}
	}

	majDiag := ""
	for i := 0; i < SIZE; i += DIM + 1 {
		majDiag += string(p.Board[i])
	}
	if isMatch(majDiag) {
		return true
	}

	minDiag := ""
	for i := DIM - 1; i <= SIZE-DIM; i += DIM - 1 {
		minDiag += string(p.Board[i])
	}
	if isMatch(minDiag) {
		return true
	}

	return false
}

type cacheKey struct {
	board string
	turn  string
}

var minimaxCache = map[cacheKey]int{}

func (p *Position) cacheKey() cacheKey {
	return cacheKey{p.Board, p.Turn}
}

func (p Position) minimax() int {
	key := p.cacheKey()

	if value, ok := minimaxCache[key]; ok {
		return value
	}

	var value int

	if p.IsWinFor("x") {
		return strings.Count(p.Board, " ")
	} else if p.IsWinFor("o") {
		return -strings.Count(p.Board, " ")
	} else if strings.Count(p.Board, " ") == 0 {
		return 0
	} else {
		if p.Turn == "x" {
			value = math.MinInt
		} else {
			value = math.MaxInt
		}

		for _, idx := range p.PossibleMoves() {
			next := p.Copy()
			next = next.Move(idx)

			v := next.minimax()

			if p.Turn == "x" {
				if v > value {
					value = v
				}
			} else {
				if v < value {
					value = v
				}
			}
		}
	}

	minimaxCache[key] = value
	return value
}

func (p *Position) Copy() *Position {
	newBoard := make([]byte, len(p.Board))
	copy(newBoard, p.Board)

	return &Position{
		Board: string(newBoard),
		Turn:  p.Turn,
	}
}

func (p Position) String() string {
	return fmt.Sprintf("%s.%s", p.Turn, p.Board)
}

func (p Position) BestMove() int {
	bestIdx := -1
	var bestVal int

	if p.Turn == "x" {
		bestVal = math.MinInt
	} else {
		bestVal = math.MaxInt
	}

	for _, idx := range p.PossibleMoves() {
		next := p.Copy()
		next.Move(idx)
		val := next.minimax()

		if p.Turn == "x" {
			if val > bestVal {
				bestVal = val
				bestIdx = idx
			}
		} else {
			if val < bestVal {
				bestVal = val
				bestIdx = idx
			}
		}
	}

	return bestIdx
}

func (p Position) IsGameEnd() bool {
	return p.IsWinFor("x") || p.IsWinFor("o") || strings.Count(p.Board, " ") == 0
}
