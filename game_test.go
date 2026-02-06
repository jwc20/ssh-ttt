package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	game := initGame()

	t.Run("test init", func(t *testing.T) {
		assert.Equal(t, game.position, initPosition())
	})

	t.Run("test quit", func(t *testing.T) {
		in := userSends("q")
		assert.Equal(t, game.askForPlayer(in), -1)
	})

	t.Run("test computer", func(t *testing.T) {
		in := userSends("1")
		assert.Equal(t, game.askForPlayer(in), COMPUTER)
	})

	t.Run("test human", func(t *testing.T) {
		in := userSends("2")
		assert.Equal(t, game.askForPlayer(in), HUMAN)
	})

	t.Run("test bad input", func(t *testing.T) {
		in := userSends("b")
		assert.Equal(t, game.askForPlayer(in), -1)
	})

	t.Run("test quit move", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})

	t.Run("test valid move", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})

	t.Run("test out of bound move", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})

	t.Run("test occupied move", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})

	t.Run("test quit play", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})

	t.Run("test choose human (play)", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})

	t.Run("test move", func(t *testing.T) {
		//assert.Equal(t, game.askForMove(), -1)
	})
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}
