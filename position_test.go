package ttt_test

import (
	"strings"
	"testing"

	ttt "github.com/jwc20/ssh-ttt"
	"github.com/stretchr/testify/assert"
)

func TestPosition(t *testing.T) {
	t.Run("test init", func(t *testing.T) {
		position := ttt.InitPosition()

		assertEqual(t, position.Board, strings.Repeat(" ", 9))
		assertEqual(t, position.Turn, "x")
	})

	t.Run("test Position with arguments", func(t *testing.T) {
		position := ttt.Position{
			Turn:  "o",
			Board: "x        ",
		}

		assertEqual(t, position.Board, "x        ")
		assertEqual(t, position.Turn, "o")
	})

	t.Run("test String()", func(t *testing.T) {
		assertEqual(t, ttt.InitPosition().String(), "x.         ")
	})

	t.Run("test equal", func(t *testing.T) {
		position := ttt.InitPosition()
		assertPositionEqual(t, *position.Move(1), ttt.Position{Turn: "o", Board: " x       "})
	})

	t.Run("test possible moves", func(t *testing.T) {
		assert.Equal(t, ttt.InitPosition().PossibleMoves(), []int{0, 1, 2, 3, 4, 5, 6, 7, 8})
		assert.Equal(t, ttt.InitPosition().Move(1).PossibleMoves(), []int{0, 2, 3, 4, 5, 6, 7, 8})
	})
}

func TestIsWinFor(t *testing.T) {
	t.Run("test no win", func(t *testing.T) {
		assert.False(t, ttt.InitPosition().IsWinFor("x"))
	})

	t.Run("test row", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "xxx      "}.IsWinFor("x"))
	})

	t.Run("test col", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "o  o  o  "}.IsWinFor("o"))
	})

	t.Run("test major diagonal", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "x   x   x"}.IsWinFor("x"))
	})

	t.Run("test minor diagonal", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "  x x x  "}.IsWinFor("x"))
	})
}

func TestBestMove(t *testing.T) {
	t.Run("test x completes winning row", func(t *testing.T) {
		assert.Equal(t, ttt.Position{Board: "xx       ", Turn: "x"}.BestMove(), 2)
	})

	t.Run("test o completes winning row", func(t *testing.T) {
		assert.Equal(t, ttt.Position{Board: "oo       ", Turn: "o"}.BestMove(), 2)
	})
}

func TestIsGameEnd(t *testing.T) {
	t.Run("test not end", func(t *testing.T) {
		assert.False(t, ttt.InitPosition().IsGameEnd())
	})

	t.Run("test end, x wins", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "xxx      "}.IsGameEnd())
	})

	t.Run("test end, o wins", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "ooo      "}.IsGameEnd())
	})

	t.Run("test end, draw", func(t *testing.T) {
		assert.True(t, ttt.Position{Board: "xoxxoxoxo"}.IsGameEnd())
	})
}

func assertPositionEqual(t *testing.T, got, want ttt.Position) {
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
