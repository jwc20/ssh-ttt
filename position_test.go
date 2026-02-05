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
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
