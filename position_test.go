package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	t.Run("test init", func(t *testing.T) {
		position := initPosition()

		assertEqual(t, position.board, strings.Repeat(" ", 9))
		assertEqual(t, position.turn, "x")
	})

	t.Run("test Position with arguments", func(t *testing.T) {
		position := Position{
			turn:  "o",
			board: "x        ",
		}

		assertEqual(t, position.board, "x        ")
		assertEqual(t, position.turn, "o")
	})

	t.Run("test String()", func(t *testing.T) {
		assertEqual(t, initPosition().String(), "x.         ")
	})

	t.Run("test equal", func(t *testing.T) {
		position := initPosition()
		assertPositionEqual(t, *position.Move(1), Position{"o", " x       "})
	})

	t.Run("test possible moves", func(t *testing.T) {
		assert.Equal(t, initPosition().PossibleMoves(), []int{0, 1, 2, 3, 4, 5, 6, 7, 8})
		assert.Equal(t, initPosition().Move(1).PossibleMoves(), []int{0, 2, 3, 4, 5, 6, 7, 8})
	})
}

func TestIsWinFor(t *testing.T) {
	t.Run("test no win", func(t *testing.T) {
		assert.False(t, initPosition().isWinFor("x"))
	})

	t.Run("test row", func(t *testing.T) {
		assert.True(t, Position{board: "xxx      "}.isWinFor("x"))
	})

	t.Run("test col", func(t *testing.T) {
		assert.True(t, Position{board: "o  o  o  "}.isWinFor("o"))
	})
	t.Run("test major diagonal", func(t *testing.T) {
		assert.True(t, Position{board: "x   x   x"}.isWinFor("x"))
	})
	t.Run("test minor diagonal", func(t *testing.T) {
		assert.True(t, Position{board: "  x x x  "}.isWinFor("x"))
	})
}

func TestMinimax(t *testing.T) {
	t.Run("test x wins", func(t *testing.T) {
		assert.Equal(t, Position{board: "xxx      "}.minimax(), 6)
	})

	t.Run("test o wins", func(t *testing.T) {
		assert.Equal(t, Position{board: "ooo      "}.minimax(), -6)
	})

	t.Run("test draw", func(t *testing.T) {
		assert.Equal(t, Position{board: "xoxxoxoxo"}.minimax(), 0)
	})

	t.Run("test x wins in one", func(t *testing.T) {
		assert.Equal(t, Position{board: "xx       ", turn: "x"}.minimax(), 6)
	})

	t.Run("test o wins in one", func(t *testing.T) {
		assert.Equal(t, Position{board: "oo       ", turn: "o"}.minimax(), -6)
	})
}

/* Asserts ****************************************************************************/

func assertPositionEqual(t *testing.T, got, want Position) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
