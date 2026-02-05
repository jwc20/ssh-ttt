package main

import (
	"strings"
	"testing"
)

func TestTicTacToe(t *testing.T) {
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
		assertEqual(t, position.Move(1).String(), Position{"o", " x       "}.String())
	})
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
