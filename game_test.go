package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	game := initGame()

	t.Run("test init", func(t *testing.T) {
		assert.Equal(t, game.position, initPosition())
	})
}
